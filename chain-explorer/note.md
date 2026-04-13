```bash
# Tìm process sử dụng port 5000
lsof -ti:5000

# Kill process đó
kill -9 $(lsof -ti:5000)