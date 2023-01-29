# Go-chat

Chat application built with go(golang channels) and WebSocket

## Package

- Kafka (code included but not deploy. Production environment uses golang
  channels)
- Gin (Web framework)
- Viper (environment variables)
- MongoDB
- Uber Zap (log)
- google proto buffer (transmit message)
- Google storage (store media files like profile image , audio, ...)
  - (update profile image) => multipart/form data will be used to transmit image
    data to the backend
  - (send media message from front-end) => send a signed URL to the front-end
    and let the front-end upload the media.

## Protocol Buffer schema

```
    string avatar = 1;
    string fromUsername = 2;
    string from = 3;
    string to = 4;
    string content = 5;
    int32 contentType = 6; // 1-> text , 2 -> pdf files etc , 3 -> pic , 4 -> audio , 5 -> video , 6-> voice call , 7 -> video call
    string type = 7;
    int32 messageType = 8; // 1 -> private message , 2- > group message
    string url = 9; // url for video , picture or audio
    string fileSuffix = 10 ;
    bytes file = 11;
```
