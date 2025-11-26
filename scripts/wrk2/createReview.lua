wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"

function init(args)
    math.randomseed(os.time())
end

function request()
    local body = string.format([[
    {
        "booking_id": %d,
        "rating": %d,
        "comment": "Test review %d"
    }
    ]], math.random(1, 100000), math.random(1, 5), math.random(1000))

    return wrk.format("POST", "/v1/review", nil, body)
end