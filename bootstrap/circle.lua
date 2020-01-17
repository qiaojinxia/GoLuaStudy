local _M = {
    rate = 3.14,
    lang = 1
}
local mt = {__index = _M}

function _M:new()
    setmetatable(_M , mt)
    return _M
end

function _M:area(lang)
    self.lang = lang
    local area = self.rate * (lang/2) * (lang/2)
    return area
end
return _M
