package main

import (
	"log"

	"github.com/AVVKavvk/ram/algo"
	socketio "github.com/googollee/go-socket.io"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Create server with proper CORS configuration
	socket := socketio.NewServer(nil)

	// Socket Connection
	socket.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("✅ Connected:", s.ID())
		return nil
	})

	// Socket Disconnection
	socket.OnDisconnect("/", func(c socketio.Conn, s string) {
		log.Println("❌ Disconnected:", c.ID())
	})

	// Socket Error
	socket.OnError("/", func(c socketio.Conn, err error) {
		log.Printf("⚠️  Error: %v, SocketId: %v", err, c.ID())
	})

	// Join Room Handler
	socket.OnEvent("/", "join", func(s socketio.Conn, room string) {
		s.Join(room)
		s.Emit("joined", room)
		log.Printf("👤 User %s joined room %s", s.ID(), room)
	})

	// Chat Message Handler
	socket.OnEvent("/", "chat", func(s socketio.Conn, msg string) {
		rooms := s.Rooms()
		log.Printf("💬 Message from %s: %s", s.ID(), msg)
		for _, room := range rooms {
			if room != s.ID() { // Don't broadcast to the connection ID room
				log.Printf("📤 Broadcasting to room: %s", room)

				// Check if message is abuse
				newNotAbuseMsg := algo.CheckAbuseAndGetNewMessage(msg)
				socket.BroadcastToRoom("/", room, "message", map[string]interface{}{"message": newNotAbuseMsg, "user": s.ID()})
			}
		}
		// Also emit back to sender
		// s.Emit("message", msg)
	})

	go func() {
		if err := socket.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer socket.Close()

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST},
	}))

	// Route Socket.IO traffic through Echo
	e.Any("/socket.io/*", func(c echo.Context) error {
		socket.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	// Standard Logger Middleware
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	log.Println("🚀 Server starting on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}

func init() {
	algo.InitTrieWithAbuseWords()
}
