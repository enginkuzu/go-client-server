package main

const ClientAppNameAndVersion = "GO TCP/IP Client v0.0.1"

const LogFileName = "client_log.txt"

var clientId uint32 = 0
var diskToSocketChannelQueue = make(chan []byte, 128)
