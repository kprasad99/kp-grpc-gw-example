package rest

import (
	"context"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/protobuf/ptypes/empty"
	service "github.com/kprasad99/kp-grpc-gateway/api/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var hopHeaders = []string{
	"Connection",          // Connection
	"Proxy-Connection",    // non-standard but still sent by libcurl and rejected by e.g. google
	"Keep-Alive",          // Keep-Alive
	"Proxy-Authenticate",  // Proxy-Authenticate
	"Proxy-Authorization", // Proxy-Authorization
	"Te",                  // canonicalized version of "TE"
	"Trailer",             // not Trailers per URL above; https://www.rfc-editor.org/errata_search.php?eid=4522
	"Transfer-Encoding",   // Transfer-Encoding
	"Upgrade",             // Upgrade
	"traceparent",
	"Content-Type",
	"Content-Length",
	"Accept",
	"Accept-Encoding",
}

func getHeaders(ctx *fiber.Ctx) map[string]string {
	headers := map[string]string{}
	found := false
	for k, v := range ctx.GetReqHeaders() {
		found = false
		for _, h := range hopHeaders {
			if strings.EqualFold(k, h) {
				found = true
				break
			}
		}
		if !found {
			headers[k] = v
		}
	}
	if clientIP, _, err := net.SplitHostPort(ctx.Context().RemoteAddr().String()); err == nil {
		hosts := ctx.Request().Header.Peek("X-Forwarded-For")
		var host string
		if hosts != nil {
			host = string(hosts) + ", " + clientIP
		}
		headers["X-Forwarded-For"] = host
	}
	return headers
}

type UserService interface {
	List(ctx *fiber.Ctx) error
	Add(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Remove(ctx *fiber.Ctx) error
}

type userService struct {
	conn *grpc.ClientConn
}

func NewUserService(cc *grpc.ClientConn) UserService {
	return &userService{cc}
}

func mapToHttpError(err error) *fiber.Error {
	st := status.Convert(err)
	log.Println("Message", st.Message())
	log.Println("Detials", st.Details())
	log.Println("Code", st.Code())
	code := 500
	switch st.Code() {
	case codes.NotFound:
		code = 404
	case codes.AlreadyExists:
		fallthrough
	case codes.InvalidArgument:
		code = 400
	default:
		code = 500
	}
	return fiber.NewError(code, st.Message())
}

func (c *userService) Add(f *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(f.UserContext(), time.Second)
	defer cancel()
	in := new(service.User)
	if err := f.BodyParser(in); err != nil {
		return err
	}
	client := service.NewUserServiceClient(c.conn)
	md := metadata.New(getHeaders(f))
	ctx = metadata.NewOutgoingContext(ctx, md)
	var respHeader metadata.MD
	out, err := client.Add(ctx, in, grpc.Header(&respHeader))
	for k, v := range respHeader {
		f.Response().Header.Set(k, strings.Join(v, ";"))
	}
	if err != nil {
		return mapToHttpError(err)
	} else {
		return f.JSON(out)
	}
}

func (c *userService) Update(f *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(f.UserContext(), time.Second)
	defer cancel()
	in := new(service.User)
	_, err := f.ParamsInt("id")
	if err != nil {
		return err
	}
	if err := f.BodyParser(in); err != nil {
		return err
	}
	client := service.NewUserServiceClient(c.conn)
	md := metadata.New(getHeaders(f))
	ctx = metadata.NewOutgoingContext(ctx, md)
	var respHeader metadata.MD
	_, err = client.Update(ctx, in, grpc.Header(&respHeader))
	for k, v := range respHeader {
		f.Response().Header.Set(k, strings.Join(v, ";"))
	}
	if err != nil {
		return mapToHttpError(err)
	}
	return f.SendStatus(200)
}

func (c *userService) Remove(f *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(f.UserContext(), time.Second)
	defer cancel()
	in := new(service.UserIDInput)
	id, err := f.ParamsInt("id")
	if err != nil {
		return err
	}
	in.Id = int32(id)
	client := service.NewUserServiceClient(c.conn)
	md := metadata.New(getHeaders(f))
	ctx = metadata.NewOutgoingContext(ctx, md)
	var respHeader metadata.MD
	_, err = client.Remove(ctx, in, grpc.Header(&respHeader))
	for k, v := range respHeader {
		f.Response().Header.Set(k, strings.Join(v, ";"))
	}
	if err != nil {
		return mapToHttpError(err)
	}
	return f.SendStatus(200)
}

func (c *userService) List(f *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(f.UserContext(), time.Second)
	defer cancel()
	in := new(empty.Empty)
	client := service.NewUserServiceClient(c.conn)
	md := metadata.New(getHeaders(f))
	ctx = metadata.NewOutgoingContext(ctx, md)
	var respHeader metadata.MD

	resp, err := client.List(ctx, in, grpc.Header(&respHeader))
	for k, v := range respHeader {
		f.Response().Header.Set(k, strings.Join(v, ";"))
	}
	if err != nil {
		return mapToHttpError(err)
	}

	users := make([]*service.User, 0)
	for {
		r, err := resp.Recv()
		if err == io.EOF {
			break
		} else if err == nil {
			users = append(users, r)
		}
		if err != nil {
			log.Println("170", err)
			return err
		}

	}
	return f.JSON(users)
}
