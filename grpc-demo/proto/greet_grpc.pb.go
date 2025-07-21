// 由 protoc-gen-go-grpc 生成的代码。不要编辑。
// 版本信息:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// 源文件: proto/greet.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// 编译时断言，确保生成的文件与当前 grpc 包兼容。
// 需要 gRPC-Go v1.64.0 或更高版本。
const _ = grpc.SupportPackageIsVersion9

const (
	// 定义服务方法的完整 gRPC 路径
	Greeter_SayHello_FullMethodName = "/greet.Greeter/SayHello"
)

// GreeterClient 是 Greeter 服务的客户端接口。
//
// 有关 ctx 使用和关闭/结束流式 RPC 的语义，请参考 https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream。
type GreeterClient interface {
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
}

// greeterClient 实现了 GreeterClient 接口
type greeterClient struct {
	cc grpc.ClientConnInterface
}

// NewGreeterClient 创建一个新的 Greeter 客户端实例
func NewGreeterClient(cc grpc.ClientConnInterface) GreeterClient {
	return &greeterClient{cc}
}

// SayHello 实现 GreeterClient 接口的方法
func (c *greeterClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, Greeter_SayHello_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GreeterServer 是 Greeter 服务的服务端接口。
// 所有实现必须嵌入 UnimplementedGreeterServer 以保持向前兼容性。
type GreeterServer interface {
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
	mustEmbedUnimplementedGreeterServer()
}

// UnimplementedGreeterServer 必须被嵌入以保证实现的向前兼容。
//
// 注意：为了防止方法调用时发生空指针解引用，应该以值类型而非指针类型嵌入。
type UnimplementedGreeterServer struct{}

func (UnimplementedGreeterServer) SayHello(context.Context, *HelloRequest) (*HelloReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}

func (UnimplementedGreeterServer) mustEmbedUnimplementedGreeterServer() {}
func (UnimplementedGreeterServer) testEmbeddedByValue()                 {}

// UnsafeGreeterServer 可以被嵌入以放弃此服务的向前兼容性。
// 不建议使用此接口，因为向 GreeterServer 添加方法将导致编译错误。
type UnsafeGreeterServer interface {
	mustEmbedUnimplementedGreeterServer()
}

// RegisterGreeterServer 注册 GreeterServer 到 gRPC 服务注册器中
func RegisterGreeterServer(s grpc.ServiceRegistrar, srv GreeterServer) {
	// 如果下面的调用 panic，说明 UnimplementedGreeterServer 是以指针形式嵌入且为 nil。
	// 这会导致在运行时调用未实现的方法时发生 panic，
	// 因此我们在初始化时进行检查以避免这种情况。
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Greeter_ServiceDesc, srv)
}

// _Greeter_SayHello_Handler 是 SayHello 方法的 gRPC 处理程序
func _Greeter_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Greeter_SayHello_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Greeter_ServiceDesc 是 Greeter 服务的 gRPC 服务描述符。
// 主要用于直接与 grpc.RegisterService 配合使用，
// 不应被内省或修改（即使是副本）
var Greeter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "greet.Greeter",
	HandlerType: (*GreeterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Greeter_SayHello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/greet.proto",
}
