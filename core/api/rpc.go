package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"localp2p/discovery"
	"localp2p/transport"
)

type RPCServer struct {
	discovery *discovery.Discovery
	transport *transport.Transport
	port      int
}

type RPCRequest struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type RPCResponse struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func NewRPCServer(disc *discovery.Discovery, trans *transport.Transport, port int) *RPCServer {
	return &RPCServer{
		discovery: disc,
		transport: trans,
		port:      port,
	}
}

func (s *RPCServer) Start() error {
	http.HandleFunc("/rpc", s.handleRPC)
	log.Printf("RPC server starting on port %d", s.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

func (s *RPCServer) handleRPC(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, "Invalid JSON request")
		return
	}

	var response RPCResponse

	switch req.Method {
	case "discover":
		peers := s.discovery.GetPeers()
		response.Result = peers
	case "connect":
		params, ok := req.Params.(map[string]interface{})
		if !ok {
			response.Error = "Invalid parameters"
		} else {
			address, _ := params["address"].(string)
			port, _ := params["port"].(float64)
			err := s.transport.ConnectToPeer(address, int(port))
			if err != nil {
				response.Error = err.Error()
			} else {
				response.Result = "Connected successfully"
			}
		}
	case "send":
		params, ok := req.Params.(map[string]interface{})
		if !ok {
			response.Error = "Invalid parameters"
		} else {
			to, _ := params["to"].(string)
			content, _ := params["content"].(string)
			err := s.transport.SendMessage(to, content)
			if err != nil {
				response.Error = err.Error()
			} else {
				response.Result = "Message sent"
			}
		}
	case "connections":
		connections := s.transport.GetConnections()
		response.Result = connections
	default:
		response.Error = "Unknown method"
	}

	json.NewEncoder(w).Encode(response)
}

func (s *RPCServer) sendError(w http.ResponseWriter, message string) {
	response := RPCResponse{Error: message}
	json.NewEncoder(w).Encode(response)
}
