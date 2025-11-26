wrk.method = "DELETE"

function init(args)
    math.randomseed(os.time())
end

function request()
    local id = math.random(1, 1000000)
    local path = "/v1/booking/" .. id
    return wrk.format("DELETE", path)
end