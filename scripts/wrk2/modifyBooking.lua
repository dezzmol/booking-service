wrk.method = "PUT"
wrk.headers["Content-Type"] = "application/json"

function init(args)
    math.randomseed(os.time())
end

function request()
    local id = math.random(1, 1000000)
    local path = "/v1/booking/" .. id

    local body = string.format([[
    {
        "new_room_id": %d,
        "new_from": "2025-02-%02d",
        "new_to": "2025-02-%02d"
    }
    ]], math.random(1, 50000), math.random(1, 27), math.random(2, 28))

    return wrk.format("PUT", path, nil, body)
end