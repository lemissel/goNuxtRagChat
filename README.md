# How to install ollama models (Install ollama locally before)
ollama pull mxbai-embed-large  
ollama pull llama2

# Run Ollama
ollama serve

# Run backend
docker-compose up
cd backend
go run .

# Run frontend
npm install
npm run dev
