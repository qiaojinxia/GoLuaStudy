local tableA = {k2 = "Hello"}
local tableB = {k1 = "Hi"}
setmetatable(tableB,{__index = tableA})
print(tableB.k2)