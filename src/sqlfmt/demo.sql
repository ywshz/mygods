  select a,b,c from tablea
          left join ( select a,c,b from tableb ) b
          on tablea.a=b.a
        where tablea.a=123 and b.b in (2, 3 ,1)