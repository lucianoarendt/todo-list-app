duas:
-tira "%s" e regex
-"%s" só ids lista

Create       Done                 Done
Read         Done                 Done
ReadAll      TODO validar cache   Done
Update       Done                 Done
Delete       Done                 Done
CreateSymbol Done                 Done
DeleteSymbol Done                 Done

Funcao      Lista     Usuario
Create      OK        NotOK/ok
Read        OK        NotOK/ok
ReadAll     OK        OK
Update      OK        NotOK/ok
CreateSym   NotOK/ok  NotOK/ok
DeleteSym   NotOK/ok  NotOK/ok

Create 1-> Read 2(OK) -> Update 1,2-> Read 1,2(OK)-> CreateSymbol 1,2 -> Read 1,2 (OK)-> DeleteSymbol 1,2 (OK)

 -> CreateSymbol, DeleteSymbol ? (NoCache)

ReadAll -> all (OK)

~~~
-- middleware
-- Refatorar ReadAllCache, Update, CreateSymbol

Update isDefault editavel?
DeleteFromCache (List ou userID e id)

Create, Update and Delete default list in cache

Refactor DDD:
    CreateSymbolOnCache listID on symbol param return err if nil
    Symbols Delete user validator