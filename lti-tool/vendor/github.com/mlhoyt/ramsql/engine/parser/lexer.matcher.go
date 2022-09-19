package parser

func (l *lexer) MatchActionToken() bool {
  return l.Match([]byte("action"), ActionToken)
}

func (l *lexer) MatchAndToken() bool {
  return l.Match([]byte("and"), AndToken)
}

func (l *lexer) MatchAsToken() bool {
  return l.Match([]byte("as"), AsToken)
}

func (l *lexer) MatchAscToken() bool {
  return l.Match([]byte("asc"), AscToken)
}

func (l *lexer) MatchAutoincrementToken() bool {
  return l.Match([]byte("autoincrement"), AutoincrementToken) ||
    l.Match([]byte("auto_increment"), AutoincrementToken)
}

func (l *lexer) MatchBacktickToken() bool {
  return l.MatchSingle('`', BacktickToken)
}

func (l *lexer) MatchBracketClosingToken() bool {
  return l.MatchSingle(')', BracketClosingToken)
}

func (l *lexer) MatchBracketOpeningToken() bool {
  return l.MatchSingle('(', BracketOpeningToken)
}

func (l *lexer) MatchBtreeToken() bool {
  return l.Match([]byte("btree"), BtreeToken)
}

func (l *lexer) MatchByToken() bool {
  return l.Match([]byte("by"), ByToken)
}

func (l *lexer) MatchCascadeToken() bool {
  return l.Match([]byte("cascade"), CascadeToken)
}

func (l *lexer) MatchCharacterToken() bool {
  return l.Match([]byte("character"), CharacterToken)
}

func (l *lexer) MatchCharsetToken() bool {
  return l.Match([]byte("charset"), CharsetToken)
}

func (l *lexer) MatchCommaToken() bool {
  return l.MatchSingle(',', CommaToken)
}

func (l *lexer) MatchConstraintToken() bool {
  return l.Match([]byte("constraint"), ConstraintToken)
}

func (l *lexer) MatchCountToken() bool {
  return l.Match([]byte("count"), CountToken)
}

func (l *lexer) MatchCreateToken() bool {
  return l.Match([]byte("create"), CreateToken)
}

func (l *lexer) MatchDefaultToken() bool {
  return l.Match([]byte("default"), DefaultToken)
}

func (l *lexer) MatchDeleteToken() bool {
  return l.Match([]byte("delete"), DeleteToken)
}

func (l *lexer) MatchDescToken() bool {
  return l.Match([]byte("desc"), DescToken)
}

func (l *lexer) MatchDropToken() bool {
  return l.Match([]byte("drop"), DropToken)
}

func (l *lexer) MatchEngineToken() bool {
  return l.Match([]byte("engine"), EngineToken)
}

func (l *lexer) MatchEqualityToken() bool {
  return l.MatchSingle('=', EqualityToken)
}

func (l *lexer) MatchExistsToken() bool {
  return l.Match([]byte("exists"), ExistsToken)
}

func (l *lexer) MatchFalseToken() bool {
  return l.Match([]byte("false"), FalseToken)
}

func (l *lexer) MatchForToken() bool {
  return l.Match([]byte("for"), ForToken)
}

func (l *lexer) MatchForeignToken() bool {
  return l.Match([]byte("foreign"), ForeignToken)
}

func (l *lexer) MatchFromToken() bool {
  return l.Match([]byte("from"), FromToken)
}

func (l *lexer) MatchFullToken() bool {
  return l.Match([]byte("full"), FullToken)
}

func (l *lexer) MatchGrantToken() bool {
  return l.Match([]byte("grant"), GrantToken)
}

func (l *lexer) MatchGreaterOrEqualToken() bool {
  return l.Match([]byte(">="), GreaterOrEqualToken)
}

func (l *lexer) MatchHashToken() bool {
  return l.Match([]byte("hash"), HashToken)
}

func (l *lexer) MatchIfToken() bool {
  return l.Match([]byte("if"), IfToken)
}

func (l *lexer) MatchInToken() bool {
  return l.Match([]byte("in"), InToken)
}

func (l *lexer) MatchIndexToken() bool {
  return l.Match([]byte("index"), IndexToken)
}

func (l *lexer) MatchInnerToken() bool {
  return l.Match([]byte("inner"), InnerToken)
}

func (l *lexer) MatchInsertToken() bool {
  return l.Match([]byte("insert"), InsertToken)
}

func (l *lexer) MatchIntoToken() bool {
  return l.Match([]byte("into"), IntoToken)
}

func (l *lexer) MatchIsToken() bool {
  return l.Match([]byte("is"), IsToken)
}

func (l *lexer) MatchJoinToken() bool {
  return l.Match([]byte("join"), JoinToken)
}

func (l *lexer) MatchKeyToken() bool {
  return l.Match([]byte("key"), KeyToken)
}

func (l *lexer) MatchLeftToken() bool {
  return l.Match([]byte("left"), LeftToken)
}

func (l *lexer) MatchLeftDipleToken() bool {
  return l.MatchSingle('<', LeftDipleToken)
}

func (l *lexer) MatchLessOrEqualToken() bool {
  return l.Match([]byte("<="), LessOrEqualToken)
}

func (l *lexer) MatchLimitToken() bool {
  return l.Match([]byte("limit"), LimitToken)
}

func (l *lexer) MatchLocalTimestampToken() bool {
  return l.Match([]byte("localtimestamp"), LocalTimestampToken) ||
    l.Match([]byte("current_timestamp"), LocalTimestampToken)
}

func (l *lexer) MatchMatchToken() bool {
  return l.Match([]byte("match"), MatchToken)
}

func (l *lexer) MatchNoToken() bool {
  return l.Match([]byte("no"), NoToken)
}

func (l *lexer) MatchNotToken() bool {
  return l.Match([]byte("not"), NotToken)
}

func (l *lexer) MatchNowToken() bool {
  return l.Match([]byte("now()"), NowToken)
}

func (l *lexer) MatchNullToken() bool {
  return l.Match([]byte("null"), NullToken)
}

func (l *lexer) MatchOffsetToken() bool {
  return l.Match([]byte("offset"), OffsetToken)
}

func (l *lexer) MatchOnToken() bool {
  return l.Match([]byte("on"), OnToken)
}

func (l *lexer) MatchOrToken() bool {
  return l.Match([]byte("or"), OrToken)
}

func (l *lexer) MatchOrderToken() bool {
  return l.Match([]byte("order"), OrderToken)
}

func (l *lexer) MatchOuterToken() bool {
  return l.Match([]byte("outer"), OuterToken)
}

func (l *lexer) MatchPartialToken() bool {
  return l.Match([]byte("partial"), PartialToken)
}

func (l *lexer) MatchPeriodToken() bool {
  return l.MatchSingle('.', PeriodToken)
}

func (l *lexer) MatchPrimaryToken() bool {
  return l.Match([]byte("primary"), PrimaryToken)
}

func (l *lexer) MatchReferencesToken() bool {
  return l.Match([]byte("references"), ReferencesToken)
}

func (l *lexer) MatchRestrictToken() bool {
  return l.Match([]byte("restrict"), RestrictToken)
}

func (l *lexer) MatchReturningToken() bool {
  return l.Match([]byte("returning"), ReturningToken)
}

func (l *lexer) MatchRightToken() bool {
  return l.Match([]byte("right"), RightToken)
}

func (l *lexer) MatchRightDipleToken() bool {
  return l.MatchSingle('>', RightDipleToken)
}

func (l *lexer) MatchSelectToken() bool {
  return l.Match([]byte("select"), SelectToken)
}

func (l *lexer) MatchSemicolonToken() bool {
  return l.MatchSingle(';', SemicolonToken)
}

func (l *lexer) MatchSetToken() bool {
  return l.Match([]byte("set"), SetToken)
}

func (l *lexer) MatchSimpleToken() bool {
  return l.Match([]byte("simple"), SimpleToken)
}

func (l *lexer) MatchStarToken() bool {
  return l.MatchSingle('*', StarToken)
}

func (l *lexer) MatchTableToken() bool {
  return l.Match([]byte("table"), TableToken)
}

func (l *lexer) MatchTimeToken() bool {
  return l.Match([]byte("time"), TimeToken)
}

func (l *lexer) MatchTrueToken() bool {
  return l.Match([]byte("true"), TrueToken)
}

func (l *lexer) MatchTruncateToken() bool {
  return l.Match([]byte("truncate"), TruncateToken)
}

func (l *lexer) MatchUniqueToken() bool {
  return l.Match([]byte("unique"), UniqueToken)
}

func (l *lexer) MatchUpdateToken() bool {
  return l.Match([]byte("update"), UpdateToken)
}

func (l *lexer) MatchUsingToken() bool {
  return l.Match([]byte("using"), UsingToken)
}

func (l *lexer) MatchValuesToken() bool {
  return l.Match([]byte("values"), ValuesToken)
}

func (l *lexer) MatchWhereToken() bool {
  return l.Match([]byte("where"), WhereToken)
}

func (l *lexer) MatchWithToken() bool {
  return l.Match([]byte("with"), WithToken)
}

func (l *lexer) MatchZoneToken() bool {
  return l.Match([]byte("zone"), ZoneToken)
}
