import io from "socket.io-client";
import readline from "readline";

const PORT = 8080;
const socket = io(`http://localhost:${PORT}`, {
  transports: ["websocket", "polling"], // Try websocket first
});

const room = "chat_room";

socket.on("connect", () => {
  console.log("✅ Connected! Socket ID:", socket.id);
  console.log("Joining room:", room);
  socket.emit("join", room);
});

socket.on("joined", (r: any) => {
  console.log("✅ Successfully joined room:", r);
  console.log("\n📝 Type your messages and press Enter to send:");
  console.log("-------------------------------------------");
  startCLI();
});

socket.on("message", (msg: any) => {
  console.log("📨 Received:", msg);
});

socket.on("disconnect", () => {
  console.log("❌ Disconnected!");
  process.exit(0);
});

socket.on("connect_error", (err: any) => {
  console.error("Connection error:", err.message);
});

// Interactive CLI for sending messages
function startCLI() {
  const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout,
    prompt: "> ",
  });

  rl.prompt();

  rl.on("line", (line) => {
    const msg = line.trim();
    if (msg) {
      if (msg.toLowerCase() === "exit" || msg.toLowerCase() === "quit") {
        console.log("👋 Goodbye!");
        socket.disconnect();
        rl.close();
        process.exit(0);
      } else {
        socket.emit("chat", msg);
        console.log("📤 Sent:", msg);
      }
    }
    rl.prompt();
  });

  rl.on("close", () => {
    socket.disconnect();
    process.exit(0);
  });
}
