wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"

function init(args)
    math.randomseed(os.time())
end

function request()
    local body = string.format([[
    {
        "name": "Guest %d",
        "email": "guest%d@example.com"
    }
    ]], math.random(100000), math.random(100000))

    return wrk.format("POST", "/v1/guests", nil, body)
end