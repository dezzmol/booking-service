wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"

function init(args)
    math.randomseed(os.time())
end

function request()
    local body = string.format([[
    {
        "name": "Hotel %d",
        "location": "City %d",
        "stars": %d
    }
    ]], math.random(100000), math.random(100), math.random(1,5))

    return wrk.format("POST", "/v1/hotels", nil, body)
end