package client
// 负责把rpc技术和领域模型分离，领域模型不能又具体的技术实现逻辑

// 负责把domainservice向外提供能力，服务之间交互不能通过包引用，而是通过client进行数据交流

// 服务主要是操作infrastructure层的逻辑，而不能直接操作internal里面模型，为了以后internal的领域模型可以单独抽出来，对上层逻辑透明


//



