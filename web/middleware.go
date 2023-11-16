package web

// MiddleWare 函数式的责任链模式 , 函数式的洋葱模式
type MiddleWare func(next HandleFunc) HandleFunc
