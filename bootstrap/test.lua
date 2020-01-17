local sum = 0
local a ={"caomao","caomaod"}
print(a[1])
for i = 100, 0 ,-1  do
  if i % 2 == 0 then
    sum = sum + i
    print(sum)
  end
end