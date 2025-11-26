wrk.method = "PUT"
wrk.headers["Content-Type"] = "application/json"

function init(args)
    math.randomseed(os.time())
end

function request()
    local body = string.format([[
    {
        "room_id": %d,
        "capacity": %d,
        "price": %d
    }
    ]], math.random(1, 50000), math.random(1, 6), math.random(50, 500))

    return wrk.format("PUT", "/v1/room", nil, body)
end