# Owl Parser

This parser combines tokens into a more useful abstract syntax tree (AST). There are a number of ways to do this, but the Owl parser currently uses a mix of LL(1) for statements (since they all start with a distinct token) and Pratt parsing for expressions, since it allows for a simple model to implement infix operators. The language grammar is the following

```
Program -> Block

Block -> Statement*

Statement -> Let 
           | Assign 
           | For 
           | While 
           | If 
           | Throw 
           | Try 
           | Expression 
           | Return 
           | Break 
           | Continue

Let -> <LET> Assignment <ASSIGN> Expression
Assign -> Assignment <ASSIGN> Expression
For -> <FOR> Assignment <IN> Expression <LBRACE> Block <RBRACE>
While -> <WHILE> Expression <LBRACE> Block <RBRACE>
If -> <IF> Expression <LBRACE> Block <RBRACE> (<ELSE> <LBRACE> Block <RBRACE>)?
Throw -> <THROW> Expression
Try -> <TRY> Block <LBRACE> Block <RBRACE> (<CATCH> <LBRACE> Block <RBRACE>)? (<FINALLY> <LBRACE> Block <RBRACE>)?
ExpressionStmt -> Expression
Return -> <RETURN> Expression
Break -> <BREAK>
Continue -> <CONTINUE>


Expression = BinOp
           | UnaryOp
           | FunctionCall
           | IfExpression
           | Map
           | Set
           | List
           | Constant
           | Function
           | Attribute
           | Index
           | Name

BinOp = Expression <PLUS> Expression
      | Expression <MINUS> Expression
      | Expression <STAR> Expression
      | Expression <SLASH> Expression
      | Expression <COMPARE> Expression
      | Expression <BOOLCOMP> Expression
      | Expression <PERCENT> Expression

UnaryOp = <MINUS> Expression
        | <BANG> Expression

FunctionCall = Expression <LPAREN> Arguments <RPAREN>
Arguments = Expression (<COMMA> Expression)*

IfExpression = Expression <QUESTION> Expression <COLON> Expression

Map = <LBRACE> MapTerms <RBRACE>
MapTerms = Expression <COLON> Expression (<COMMA> Expression <COLON> Expression)*

Set = <LBRACKET> SetTerms <RBRACKET>
SetTerms = Expression (<COMMA> Expression)*

List = <LBRACKET> ListTerms <RBRACKET>
ListTerms = Expression (<COMMA> Expression)*

Constant = <NUMBER>
         | <BOOL>
         | <STRING>
         | <NULL>

Function = <LParen> Assignment <RPAREN> <ARROW> Expression
         | <LParen> Assignment <RPAREN> <ARROW> <LBRACE> Block <RBRACE>

Atrtriute = Expression <DOT> <NAME>

Index = Expression <LBRACKET> Expression <RBRACKET>

Name = <NAME>



Assignment -> (<COMMA> AssignmentTerm)*
AssignmentTerm -> NameAssign
                | AttributeAssign // Not allowed in argument lists
                | IndexAssign // Not allowed in argument lists

NameAssign = <NAME>
AttributeAssign = AssignmentTerm <DOT> <NAME>
IndexAssign = AssignmentTerm <LBRACKET> Expression <RBRACKET>


<LET> = "let"
<ASSIGN> = "="
<FOR> = "for"
<IN> = "in"
<WHILE> = "while"
<IF> = "if"
<ELSE> = "else"
<RETURN> = "return"
<THROW> = "throw"
<TRY> = "try"
<CATCH> = "catch"
<FINALLY> = "finally"
<BREAK> = "break"
<CONTINUE> = "continue"

<PLUS> = "+"
<MINUS> = "-"
<STAR> = "*"
<PERCENT> = "%"
<SLASH> = "/"
<INCDEC> = "++" | "--"
<SPREAD> = "..."

<COMPARE> = "<" | ">" | "<=" | ">=" | "==" | "!="
<BOOLCOMP> = "and" | "or"
<QUESTION> = "?"
<COLON> = ":"
<COMMA> = ","
<DOT> = "."

<LBRACE> = "{"
<RBRACE> = "}"
<LBRACKET> = "["
<RBRACKET> = "]"
<LPAREN> = "("
<RPAREN> = ")"

<NAME> = /[a-zA-Z_][a-zA-Z0-9_]*/
<NUMBER> = /[0-9]*\.?[0-9]+([eE][-+]?[0-9]+)?/
<BOOL> = "true" | "false"
<STRING> = /\".*\"/
<NULL> = "null"

<COMMENT> = "//"
<NEWLINE> = "\n"
<EOF> = /$/
```