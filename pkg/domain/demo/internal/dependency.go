package internal
// 作为领域模型的接口层，只定义接口，用于领域domain的依赖倒置于基础设施技术细节

// 如果把领域下沉了，该层就是领域模型与领域模型的share kernal的核心

// 如果在服务里面。领域只能以服务形式向外提供能力，其它服务不能直接引用领域模型，做到隔离


// 这一层对外提供的是领域实体 Entity

// 这一层依赖关系：不能与外部有任何依赖，只能通过依赖倒置依赖于基础设施
