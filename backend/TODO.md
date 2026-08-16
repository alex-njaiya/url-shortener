# feature #1: Hashing a long url to a short url
1. Take a long url, pass it onto a hashing function such md5,crc32. The short-code should be 6 characters long.
2. After the short-code is generated check whether it exists on the db
3. If the short code exists then there is a collision
4. To fight the collision we append a predefined string on the short-code 
5. If the short-code is not in the db then there is no collision hence save it and return