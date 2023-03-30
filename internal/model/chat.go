/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 14:11:49
 * @LastEditTime: 2023-03-29 14:13:04
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/internal/model/chat.go
 */
package model

type ChatEvent struct {
	Message chan string
	// New client connections
	NewClients chan chan string

	// Closed client connections
	ClosedClients chan chan string

	// Total client connections
	TotalClients map[chan string]bool
}
