/**
 * Vash - JavaScript Template Parser, v0.7.11-3
 *
 * https://github.com/kirbysayshi/vash
 *
 * Copyright (c) 2013 Andrew Petersen
 * MIT License (LICENSE)
 */

;(function(vash){
    module['exports'] = vash; // NODEJS
})(function(exports){

    var vash = exports; // neccessary for nodejs references/*jshint strict:false, asi:true, laxcomma:true, laxbreak:true, boss:true, curly:true, node:true, browser:true, devel:true */

    // The basic tokens, defined as constants
    var  AT = 'AT'
    ,ASSIGN_OPERATOR = 'ASSIGN_OPERATOR'
    ,AT_COLON = 'AT_COLON'
    ,AT_STAR_CLOSE = 'AT_STAR_CLOSE'
    ,AT_STAR_OPEN = 'AT_STAR_OPEN'
    ,BACKSLASH = 'BACKSLASH'
    ,BRACE_CLOSE = 'BRACE_CLOSE'
    ,BRACE_OPEN = 'BRACE_OPEN'
    ,CONTENT = 'CONTENT'
    ,DOUBLE_QUOTE = 'DOUBLE_QUOTE'
    ,EMAIL = 'EMAIL'
    ,ESCAPED_QUOTE = 'ESCAPED_QUOTE'
    ,FORWARD_SLASH = 'FORWARD_SLASH'
    ,FUNCTION = 'FUNCTION'
    ,HARD_PAREN_CLOSE = 'HARD_PAREN_CLOSE'
    ,HARD_PAREN_OPEN = 'HARD_PAREN_OPEN'
    ,HTML_TAG_CLOSE = 'HTML_TAG_CLOSE'
    ,HTML_TAG_OPEN = 'HTML_TAG_OPEN'
    ,HTML_TAG_VOID_CLOSE = 'HTML_TAG_VOID_CLOSE'
    ,IDENTIFIER = 'IDENTIFIER'
    ,KEYWORD = 'KEYWORD'
    ,LOGICAL = 'LOGICAL'
    ,NEWLINE = 'NEWLINE'
    ,NUMERIC_CONTENT = 'NUMERIC_CONTENT'
    ,OPERATOR = 'OPERATOR'
    ,PAREN_CLOSE = 'PAREN_CLOSE'
    ,PAREN_OPEN = 'PAREN_OPEN'
    ,PERIOD = 'PERIOD'
    ,SINGLE_QUOTE = 'SINGLE_QUOTE'
    ,TEXT_TAG_CLOSE = 'TEXT_TAG_CLOSE'
    ,TEXT_TAG_OPEN = 'TEXT_TAG_OPEN'
    ,WHITESPACE = 'WHITESPACE';

    var PAIRS = {};

    // defined through indexing to help minification
    PAIRS[AT_STAR_OPEN] = AT_STAR_CLOSE;
    PAIRS[BRACE_OPEN] = BRACE_CLOSE;
    PAIRS[DOUBLE_QUOTE] = DOUBLE_QUOTE;
    PAIRS[HARD_PAREN_OPEN] = HARD_PAREN_CLOSE;
    PAIRS[PAREN_OPEN] = PAREN_CLOSE;
    PAIRS[SINGLE_QUOTE] = SINGLE_QUOTE;
    PAIRS[AT_COLON] = NEWLINE;
    PAIRS[FORWARD_SLASH] = FORWARD_SLASH; // regex



    // The order of these is important, as it is the order in which
    // they are run against the input string.
    // They are separated out here to allow for better minification
    // with the least amount of effort from me. :)

    // NOTE: this is an array, not an object literal! The () around
    // the regexps are for the sake of the syntax highlighter in my
    // editor... sublimetext2

    var TESTS = [

	// A real email address is considerably more complex, and unfortunately
	// this complexity makes it impossible to differentiate between an address
	// and an AT expression.
	//
	// Instead, this regex assumes the only valid characters for the user portion
	// of the address are alphanumeric, period, and %. This means that a complex email like
	// who-something@example.com will be interpreted as an email, but incompletely. `who-`
	// will be content, while `something@example.com` will be the email address.
	//
	// However, this is "Good Enough"© :).
	EMAIL, (/^([a-zA-Z0-9.%]+@[a-zA-Z0-9.\-]+\.(?:ca|co\.uk|com|edu|net|org))\b/)


	,AT_STAR_OPEN, (/^(@\*)/)
	,AT_STAR_CLOSE, (/^(\*@)/)


	,AT_COLON, (/^(@\:)/)
	,AT, (/^(@)/)


	,PAREN_OPEN, (/^(\()/)
	,PAREN_CLOSE, (/^(\))/)


	,HARD_PAREN_OPEN, (/^(\[)/)
	,HARD_PAREN_CLOSE, (/^(\])/)


	,BRACE_OPEN, (/^(\{)/)
	,BRACE_CLOSE, (/^(\})/)


	,TEXT_TAG_OPEN, (/^(<text>)/)
	,TEXT_TAG_CLOSE, (/^(<\/text>)/)


	,HTML_TAG_OPEN, function(){

	    // Some context:
	    // These only need to match something that is _possibly_ a tag,
	    // self closing tag, or email address. They do not need to be able to
	    // fully parse a tag into separate parts. They can be thought of as a
	    // huge look ahead to determine if a large swath of text is an tag,
	    // even if it contains other components (like expressions or else).

	    var  reHtml = /^(<[a-zA-Z@]+?[^>]*?["a-zA-Z]*>)/
		,reEmail = /([a-zA-Z0-9.%]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,4})\b/

	    var tok = this.scan( reHtml, HTML_TAG_OPEN );

	    if( tok ){
		this.spewIf( tok, reEmail );
		this.spewIf( tok, /(@)/ );
		this.spewIf( tok, /(\/\s*>)/ );
	    }

	    return tok;
	}
	,HTML_TAG_CLOSE, (/^(<\/[^>@\b]+?>)/)
	,HTML_TAG_VOID_CLOSE, (/^(\/\s*>)/)


	,PERIOD, (/^(\.)/)
	,NEWLINE, function(){
	    var token = this.scan(/^(\n)/, NEWLINE);
	    if(token){
		this.lineno++;
		this.charno = 0;
	    }
	    return token;
	}
	,WHITESPACE, (/^(\s)/)
	,FUNCTION, (/^(function)(?![\d\w])/)
	,KEYWORD, (/^(case|do|else|section|for|func|goto|if|return|switch|try|var|while|with)(?![\d\w])/)
	,IDENTIFIER, (/^([_$a-zA-Z\xA0-\uFFFF][_$a-zA-Z0-9\xA0-\uFFFF]*)/)

	,FORWARD_SLASH, (/^(\/)/)

	,OPERATOR, (/^(===|!==|==|!==|>>>|<<|>>|>=|<=|>|<|\+|-|\/|\*|\^|%|\:|\?)/)
	,ASSIGN_OPERATOR, (/^(\|=|\^=|&=|>>>=|>>=|<<=|-=|\+=|%=|\/=|\*=|=)/)
	,LOGICAL, (/^(&&|\|\||&|\||\^)/)


	,ESCAPED_QUOTE, (/^(\\+['"])/)
	,BACKSLASH, (/^(\\)/)
	,DOUBLE_QUOTE, (/^(\")/)
	,SINGLE_QUOTE, (/^(\')/)


	,NUMERIC_CONTENT, (/^([0-9]+)/)
	,CONTENT, (/^([^\s})@.]+?)/)

    ];

    // This pattern and basic lexer code were originally from the
    // Jade lexer, but have been modified:
    // https://github.com/visionmedia/jade/blob/master/lib/lexer.js

    function VLexer(str){
	this.input = this.originalInput = str
	    .replace(/^\uFEFF/, '') // Kill BOM
	    .replace(/\r\n|\r/g, '\n');
	this.lineno = 1;
	this.charno = 0;
    }

    VLexer.prototype = {

	scan: function(regexp, type){
	    var captures, token;
	    if (captures = regexp.exec(this.input)) {
		this.input = this.input.substr((captures[1].length));

		token = {
		    type: type
		    ,line: this.lineno
		    ,chr: this.charno
		    ,val: captures[1] || ''
		    ,toString: function(){
			return '[' + this.type
			    + ' (' + this.line + ',' + this.chr + '): '
			    + this.val + ']';
		    }
		};

		this.charno += captures[0].length;
		return token;
	    }
	}

	,spewIf: function( tok, re ){
	    var result, index, spew

	    if( tok ){
		result = re.exec( tok.val );

		if( result ){
		    index = tok.val.indexOf( result[1] );
		    spew = tok.val.substring( index );
		    this.input = spew + this.input;
		    this.charno -= spew.length;
		    tok.val = tok.val.substring( 0, index );
		}
	    }

	    return tok;
	}

	,advance: function() {

	    var i, name, test, result;

	    for(i = 0; i < TESTS.length; i += 2){
		test = TESTS[i+1];
		test.displayName = TESTS[i];

		if(typeof test === 'function'){
		    // assume complex callback
		    result = test.call(this);
		}

		if(typeof test.exec === 'function'){
		    // assume regex
		    result = this.scan(test, TESTS[i]);
		}

		if( result ){
		    return result;
		}
	    }
	}
    }
    /*jshint strict:false, asi:true, laxcomma:true, laxbreak:true, boss:true, curly:true, node:true, browser:true, devel:true */

    var vQuery = function(node){
	return new vQuery.fn.init(node);
    }

    vQuery.prototype.init = function(astNode){

	// handle mode string
	if(typeof astNode === 'string'){
	    this.mode = astNode;
	}

	this.maxCheck();
    }

    vQuery.fn = vQuery.prototype.init.prototype = vQuery.prototype;

    vQuery.fn.vquery = 'yep';
    vQuery.fn.constructor = vQuery;
    vQuery.fn.length = 0;
    vQuery.fn.parent = null;
    vQuery.fn.mode = null;
    vQuery.fn.tagName = null;

    vQuery.fn.beget = function(mode, tagName){
	var child = vQuery(mode);
	child.parent = this;
	this.push( child );

	if(tagName) { child.tagName = tagName; }

	this.maxCheck();

	return child;
    }

    vQuery.fn.closest = function(mode, tagName){
	var p = this;

	while(p){

	    if( p.tagName !== tagName && p.parent ){
		p = p.parent;
	    } else {
		break;
	    }
	}

	return p;
    }

    vQuery.fn.pushFlatten = function(node){
	var n = node, i, children;

	while( n.length === 1 && n[0].vquery ){
	    n = n[0];
	}

	if(n.mode !== PRG){
	    this.push(n);
	} else {

	    for(i = 0; i < n.length; i++){
		this.push( n[i] );
	    }
	}

	this.maxCheck();

	return this;
    }

    vQuery.fn.push = function(nodes){

	if(vQuery.isArray(nodes)){
	    if(nodes.vquery){
		nodes.forEach(function(node){ node.parent = this; }, this);
	    }

	    Array.prototype.push.apply(this, nodes);
	} else {
	    if(nodes.vquery){
		nodes.parent = this;
	    }

	    Array.prototype.push.call(this, nodes);
	}

	this.maxCheck();

	return this.length;
    }

    vQuery.fn.root = function(){
	var p = this;

	while(p && p.parent && (p = p.parent)){}

	return p;
    }

    vQuery.fn.toTreeString = function(){
	var  buffer = []
	,indent = 1;

	function visitNode(node){
	    var  children
	    ,child;

	    buffer.push( Array(indent).join(' |') + ' +' + node.mode + ' ' + ( node.tagName || '' ) );

	    indent += 1;
	    children = node.slice();
	    while( (child = children.shift()) ){

		if(child.vquery === vQuery.fn.vquery){
		    // recurse
		    visitNode(child);
		} else {
		    buffer.push( Array(indent).join(' |') + ' '
				 + (child
				    ?  child.toString().replace(/(\r|\n)/g, '')
				    : '[empty]')
			       );
		}

	    }

	    indent -= 1;
	    buffer.push( Array(indent).join(' |') + ' -' + node.mode + ' ' + ( node.tagName || '' ) );
	}

	visitNode(this);

	return buffer.join('\n');
    }

    vQuery.fn.maxCheck = function(last){
	if( this.length >= vQuery.maxSize ){
	    var e = new Error();
	    e.message = 'Maximum number of elements exceeded.\n'
		+ 'This is typically caused by an unmatched character or tag. Parse tree follows:\n'
		+ this.toTreeString();
	    e.name = 'vQueryDepthException';
	    throw e;
	}
    }

    vQuery.maxSize = 100000;

    // takes a full nested set of vqueries (e.g. an AST), and flattens them
    // into a plain array. Useful for performing queries, or manipulation,
    // without having to handle a lot of parsing state.
    vQuery.fn.flatten = function(){
	var reduced;
	return this.reduce(function flatten(all, tok, i, orig){

	    if( tok.vquery ){
		all.push( { type: 'META', val: 'START' + tok.mode, tagName: tok.tagName } );
		reduced = tok.reduce(flatten, all);
		reduced.push( { type: 'META', val: 'END' + tok.mode, tagName: tok.tagName } );
		return reduced;
	    }

	    // grab the mode from the original vquery container
	    tok.mode = orig.mode;
	    all.push( tok );

	    return all;
	}, []);
    }

    // take a flat array created via vQuery.fn.flatten, and recreate the
    // original AST.
    vQuery.reconstitute = function(arr){
	return arr.reduce(function recon(ast, tok, i, orig){

	    if( tok.type === 'META' ) {
		ast = ast.parent;
	    } else {

		if( tok.mode !== ast.mode ) {
		    ast = ast.beget(tok.mode, tok.tagName);
		}

		ast.push( tok );
	    }

	    return ast;
	}, vQuery(PRG))
    }

    vQuery.isArray = function(obj){
	return Object.prototype.toString.call(obj) == '[object Array]';
    }

    vQuery.extend = function(obj){
	var next, i, p;

	for(i = 1; i < arguments.length; i++){
	    next = arguments[i];

	    for(p in next){
		obj[p] = next[p];
	    }
	}

	return obj;
    }

    vQuery.takeMethodsFromArray = function(){
	var methods = [
	    'pop', 'push', 'reverse', 'shift', 'sort', 'splice', 'unshift',
	    'concat', 'join', 'slice', 'indexOf', 'lastIndexOf',
	    'filter', 'forEach', 'every', 'map', 'some', 'reduce', 'reduceRight'
	]

	,arr = []
	,m;

	for (var i = 0; i < methods.length; i++){
	    m = methods[i];
	    if( typeof arr[m] === 'function' ){
		if( !vQuery.fn[m] ){
		    (function(methodName){
			vQuery.fn[methodName] = function(){
			    return arr[methodName].apply(this, Array.prototype.slice.call(arguments, 0));
			}
		    })(m);
		}
	    } else {
		throw new Error('Vash requires ES5 array iteration methods, missing: ' + m);
	    }
	}

    }

    vQuery.takeMethodsFromArray(); // run on page load
    /*jshint strict:false, asi:true, laxcomma:true, laxbreak:true, boss:true, curly:true, node:true, browser:true, devel:true */

    function VParser(tokens, options){

	this.options = options || {};
	this.tokens = tokens;
	//console.log("tokens: ", tokens)
	this.ast = vQuery(PRG);
	this.prevTokens = [];

	this.inCommentLine = false;
    }

    var PRG = "PROGRAM", MKP = "MARKUP", BLK = "BLOCK", EXP = "EXPRESSION" ;

    VParser.prototype = {

	parse: function(){
	    var curr, i, len, block;

	    while( this.prevTokens.push( curr ), (curr = this.tokens.pop()) ){

		if(this.ast.mode === PRG || this.ast.mode === null){
        	    this.ast = this.ast.beget( this.options.initialMode || MKP );
		    if(this.options.initialMode === EXP){
			this.ast = this.ast.beget( EXP ); // EXP needs to know it's within to continue
		    }
		}

		//console.log("curr", curr, this.ast.mode)
		if(this.ast.mode === MKP){
		    this.handleMKP(curr);
		    continue;
		}

		if(this.ast.mode === BLK){
		    this.handleBLK(curr);
		    continue;
		}

		if(this.ast.mode === EXP){
		    this.handleEXP(curr);
		    continue;
		}
	    }

	    this.ast = this.ast.root();

	    return this.ast;
	}

	,exceptionFactory: function(e, type, tok){

	    // second param is either a token or string?

	    if(type == 'UNMATCHED'){

		e.name = "UnmatchedCharacterError";

		this.ast = this.ast.root();

		if(tok){
		    e.message = 'Unmatched ' + tok.type
		    //+ ' near: "' + context + '"'
			+ ' at line ' + tok.line
			+ ', character ' + tok.chr
			+ '. Value: ' + tok.val
			+ '\n ' + this.ast.toTreeString();
		    e.lineNumber = tok.line;
		}
	    }

	    return e;
	}

	,advanceUntilNot: function(untilNot){
	    var curr, next, tks = [];

	    while( next = this.tokens[ this.tokens.length - 1 ] ){
		if(next.type === untilNot){
		    curr = this.tokens.pop();
		    tks.push(curr);
		} else {
		    break;
		}
	    }

	    return tks;
	}

	,advanceUntilMatched: function(curr, start, end, startEscape, endEscape){
	    var  next = curr
	    ,prev = null
	    ,nstart = 0
	    ,nend = 0
	    ,tks = [];

	    // this is fairly convoluted because the start and end for single/double
	    // quotes is the same, and can also be escaped

	    while(next){

		if( next.type === start ){

		    if( (prev && prev.type !== startEscape && start !== end) || !prev ){
			nstart++;
		    } else if( start === end && prev.type !== startEscape ) {
			nend++;
		    }

		} else if( next.type === end ){
		    nend++;
		    if(prev && prev.type === endEscape){ nend--; }
		}

		tks.push(next);

		if(nstart === nend) { break; }
		prev = next;
		next = this.tokens.pop();
		if(!next) { throw this.exceptionFactory(new Error(), 'UNMATCHED', curr); }
	    }

	    return tks.reverse();
	}

	,subParse: function(curr, modeToOpen, includeDelimsInSub){
	    var  subTokens
	    ,closer
	    ,miniParse
	    ,parseOpts = vQuery.extend({}, this.options);

	    parseOpts.initialMode = modeToOpen;

	    subTokens = this.advanceUntilMatched(
		curr
		,curr.type
		,PAIRS[ curr.type ]
		,null
		,AT );

            subTokens.pop();

            closer = subTokens.shift();

            if( !includeDelimsInSub ){
		this.ast.push(curr);
	    }

	    miniParse = new VParser( subTokens, parseOpts );
            miniParse.parse();

            if( includeDelimsInSub ){
		// attach delimiters to [0] (first child), because ast is PROGRAM
		miniParse.ast[0].unshift( curr );
		miniParse.ast[0].push( closer );
	    }

	    this.ast.pushFlatten(miniParse.ast);

	    if( !includeDelimsInSub ){
		this.ast.push(closer);
	    }
	}

	,handleMKP: function(curr){
	    var  next = this.tokens[ this.tokens.length - 1 ]
	    ,ahead = this.tokens[ this.tokens.length - 2 ]
	    ,tagName = null
	    ,opener;

	    switch(curr.type){

	    case AT_STAR_OPEN:
		this.advanceUntilMatched(curr, AT_STAR_OPEN, AT_STAR_CLOSE, AT, AT);
		break;

	    case AT:
		if(next) {

		    if(this.options.saveAT)  {
			this.ast.push( curr );
		    }
		    switch(next.type){

		    case PAREN_OPEN:
		    case IDENTIFIER:

			if(this.ast.length === 0) {
			    this.ast = this.ast.parent;
			    this.ast.pop(); // remove empty MKP block
			}

			this.ast = this.ast.beget( EXP );
			break;

		    case KEYWORD:
		    case FUNCTION:
		    case BRACE_OPEN:

			if(this.ast.length === 0) {
			    this.ast = this.ast.parent;
			    this.ast.pop(); // remove empty MKP block
			}

			this.ast = this.ast.beget( BLK );
			break;

		    case AT:
		    case AT_COLON:

			// we want to keep the token, but remove its
			// "special" meaning because during compilation
			// AT and AT_COLON are discarded
			next.type = 'CONTENT';
			this.ast.push( this.tokens.pop() );
			break;

		    default:
			this.ast.push( this.tokens.pop() );
			break;
		    }

		}
		break;

	    case TEXT_TAG_OPEN:
	    case HTML_TAG_OPEN:
		tagName = curr.val.match(/^<([^\/ >]+)/i);
		if(tagName === null && next && next.type === AT && ahead){
		    tagName = ahead.val.match(/(.*)/); // HACK for <@exp>
		}

		if(this.ast.tagName){
		    // current markup is already waiting for a close tag, make new child
		    this.ast = this.ast.beget(MKP, tagName[1]);
		} else {
		    this.ast.tagName = tagName[1];
		}

		if(
		    HTML_TAG_OPEN === curr.type
			|| this.options.saveTextTag
		){
		    this.ast.push(curr);
		}

		break;

	    case TEXT_TAG_CLOSE:
	    case HTML_TAG_CLOSE:
		tagName = curr.val.match(/^<\/([^>]+)/i);

		if(tagName === null && next && next.type === AT && ahead){
		    tagName = ahead.val.match(/(.*)/); // HACK for </@exp>
		}

		opener = this.ast.closest( MKP, tagName[1] );

		if(opener === null || opener.tagName !== tagName[1]){
		    // couldn't find opening tag
		    // could mean this closer is within a child parser
		    //throw this.exceptionFactory(new Error, 'UNMATCHED', curr);
		} else {
		    this.ast = opener;
		}

		if(HTML_TAG_CLOSE === curr.type || this.options.saveTextTag) {
		    this.ast.push( curr );
		}

		// close this ast if parent is BLK. if another tag follows, BLK will
		// flip over to MKP
		// Yukang: vash.js BUG here, should flip current MKP
		// so that we keep a right hierarchy
		if( this.ast.parent && this.ast.parent.mode === BLK ){
		    this.ast = this.ast.parent;
		}

		break;

	    case HTML_TAG_VOID_CLOSE:
		this.ast.push(curr);
		this.ast = this.ast.parent;
		break;

	    case BACKSLASH:
		curr.val += '\\';
		this.ast.push(curr);
		break;

	    default:
		this.ast.push(curr);
		break;
	    }

	}

	,handleBLK: function(curr){

	    var  next = this.tokens[ this.tokens.length - 1 ]
	    ,submode
	    ,opener
	    ,closer
	    ,subTokens
	    ,parseOpts
	    ,miniParse
	    ,i;

	    switch(curr.type){

	    case AT:
		if(next.type !== AT && !this.inCommentLine){
		    this.tokens.push(curr); // defer
		    this.ast = this.ast.beget(MKP);
		} else {
		    // we want to keep the token, but remove its
		    // "special" meaning because during compilation
		    // AT and AT_COLON are discarded
		    next.type = CONTENT;
		    this.ast.push(next);
		    this.tokens.pop(); // skip following AT
		}
		break;

	    case AT_STAR_OPEN:
		this.advanceUntilMatched(curr, AT_STAR_OPEN, AT_STAR_CLOSE, AT, AT);
		break;

	    case AT_COLON:
		this.subParse(curr, MKP, true);
		break;

	    case TEXT_TAG_OPEN:
	    case TEXT_TAG_CLOSE:
	    case HTML_TAG_OPEN:
	    case HTML_TAG_CLOSE:
		this.ast = this.ast.beget(MKP);
		this.tokens.push(curr); // defer
		break;

	    case FORWARD_SLASH:
	    case SINGLE_QUOTE:
	    case DOUBLE_QUOTE:
		if(
		    curr.type === FORWARD_SLASH
			&& next
			&& next.type === FORWARD_SLASH
		){
		    this.inCommentLine = true;
		}

		if(!this.inCommentLine) {
		    // assume regex or quoted string
		    subTokens = this.advanceUntilMatched(
			curr
			,curr.type
			,PAIRS[ curr.type ]
			,BACKSLASH
			,BACKSLASH ).map(function(tok){
			    // mark AT within a regex/quoted string as literal
			    if(tok.type === AT) tok.type = CONTENT;
			    return tok;
			});
		    this.ast.pushFlatten(subTokens.reverse());
		} else {
		    this.ast.push(curr);
		}

		break;

	    case NEWLINE:
		if(this.inCommentLine){
		    this.inCommentLine = false;
		}
		this.ast.push(curr);
		break;

	    case BRACE_OPEN:
	    case PAREN_OPEN:
		submode = this.options.favorText && curr.type === BRACE_OPEN
		    ? MKP
		    : BLK;

		this.subParse( curr, submode );

		subTokens = this.advanceUntilNot(WHITESPACE);
		next = this.tokens[ this.tokens.length - 1 ];
		if(
		    next
			&& next.type !== KEYWORD
			&& next.type !== FUNCTION
			&& next.type !== BRACE_OPEN
			&& curr.type !== PAREN_OPEN
		){
		    // defer whitespace
		    this.tokens.push.apply(this.tokens, subTokens.reverse());
		    this.ast = this.ast.parent;
		} else {
		    this.ast.push(subTokens);
		}

		break;

	    default:
		this.ast.push(curr);
		break;
	    }

	}

	,handleEXP: function(curr){

	    var ahead = null
	    ,opener
	    ,closer
	    ,parseOpts
	    ,miniParse
	    ,subTokens
	    ,prev
	    ,i;

	    switch(curr.type){

	    case KEYWORD:
	    case FUNCTION:
		this.ast = this.ast.beget(BLK);
		this.tokens.push(curr); // defer
		break;

	    case WHITESPACE:
	    case LOGICAL:
	    case ASSIGN_OPERATOR:
	    case OPERATOR:
	    case NUMERIC_CONTENT:
		if(this.ast.parent && this.ast.parent.mode === EXP){

		    this.ast.push(curr);
		} else {

		    // if not contained within a parent EXP, must be end of EXP
		    this.ast = this.ast.parent;
		    this.tokens.push(curr); // defer
		}

		break;

	    case IDENTIFIER:
		this.ast.push(curr);
		break;

	    case SINGLE_QUOTE:
	    case DOUBLE_QUOTE:

		if(this.ast.parent && this.ast.parent.mode === EXP){
		    subTokens = this.advanceUntilMatched(
			curr
			,curr.type
			,PAIRS[ curr.type ]
			,BACKSLASH
			,BACKSLASH );
		    this.ast.pushFlatten(subTokens.reverse());

		} else {
		    // probably end of expression
		    this.ast = this.ast.parent;
		    this.tokens.push(curr); // defer
		}

		break;

	    case HARD_PAREN_OPEN:
	    case PAREN_OPEN:

		prev = this.prevTokens[ this.prevTokens.length - 1 ];
		ahead = this.tokens[ this.tokens.length - 1 ];

		if( curr.type === HARD_PAREN_OPEN && ahead.type === HARD_PAREN_CLOSE ){
		    // likely just [], which is not likely valid outside of EXP
		    this.tokens.push(curr); // defer
		    this.ast = this.ast.parent; //this.ast.beget(MKP);
		    break;
		}

		this.subParse(curr, EXP);
		ahead = this.tokens[ this.tokens.length - 1 ];

		if( (prev && prev.type === AT) || (ahead && ahead.type === IDENTIFIER) ){
		    // explicit expression is automatically ended
		    this.ast = this.ast.parent;
		}

		break;

	    case BRACE_OPEN:
		this.tokens.push(curr); // defer
		this.ast = this.ast.beget(BLK);
		break;

	    case PERIOD:
		ahead = this.tokens[ this.tokens.length - 1 ];
		if(
		    ahead &&
			(  ahead.type === IDENTIFIER
			   || ahead.type === KEYWORD
			   || ahead.type === FUNCTION
			   || ahead.type === PERIOD
			   // if it's "expressions all the way down", then there is no way
			   // to exit EXP mode without running out of tokens, i.e. we're
			   // within a sub parser
			   || this.ast.parent && this.ast.parent.mode === EXP )
		) {
		    this.ast.push(curr);
		} else {
		    this.ast = this.ast.parent;
		    this.tokens.push(curr); // defer
		}
		break;

	    default:

		if( this.ast.parent && this.ast.parent.mode !== EXP ){
		    // assume end of expression
		    this.ast = this.ast.parent;
		    this.tokens.push(curr); // defer
		} else {
		    this.ast.push(curr);
		}

		break;
	    }
	}
    }
    /*jshint strict:false, asi:true, laxcomma:true, laxbreak:true, boss:true, curly:true, node:true, browser:true, devel:true */

    function VCompiler(ast, originalMarkup, options){
	this.ast = ast;
	this.originalMarkup = originalMarkup || '';
	this.options = options || {};

	this.reQuote = /(['"])/gi
	this.reEscapedQuote = /\\+(["'])/gi
	this.reLineBreak = /\r?\n/gi
	this.reHelpersName = /HELPERSNAME/g
	this.reModelName = /MODELNAME/g
	this.reOriginalMarkup = /ORIGINALMARKUP/g

	this.buffer = [];
    }

    var VCP = VCompiler.prototype;

    VCP.visitMarkupTok = function(tok, parentNode, index){

	this.buffer.push(
	    "MKP(" + tok.val
		.replace(this.reEscapedQuote, '\\\\$1')
		.replace(this.reQuote, '\\$1')
		.replace(this.reLineBreak, '\\n')
		+ ")MKP" );
    }

    VCP.visitBlockTok = function(tok, parentNode, index){
	this.buffer.push( "BLK(" + tok.val + ")BLK");
    }

    VCP.visitExpressionTok = function(tok, parentNode, index, isHomogenous){

	var  start = ''
	,end = ''
	,parentParentIsNotEXP = parentNode.parent && parentNode.parent.mode !== EXP;

	if(this.options.htmlEscape !== false){
	    if( parentParentIsNotEXP && index === 0 && isHomogenous ){
		if (tok.val == 'helper' || tok.val == 'raw' || this.options['package'] == 'layout') {
		    //todo: this actually results in: _buffer.WriteString((u.Intro))
		    //      should remove the extra braket
		    start += '(';
		} else {
		    start += 'gorazor.HTMLEscape(';
		}
	    }

	    if( parentParentIsNotEXP && index === parentNode.length - 1 && isHomogenous){
		end += ")";
	    }
	}

	if(parentParentIsNotEXP && (index === 0 ) ){
	    start = "_buffer.WriteString(" + start;
	}

	if( parentParentIsNotEXP && index === parentNode.length - 1 ){
	    end += ")\n";
	}

	if (tok.val == "raw") {
	    this.buffer.push( start + end);
	} else {
	    this.buffer.push( start + tok.val + end );
	}
    }

    VCP.visitNode = function(node){

	var n, children = node.slice(0), nonExp, i, child;

	if(node.mode === EXP && (node.parent && node.parent.mode !== EXP)){
	    // see if this node's children are all EXP
	    nonExp = node.filter(VCompiler.findNonExp).length;
	}

	for(i = 0; i < children.length; i++){
	    child = children[i];

	    // if saveAT is true, or if AT_COLON is used, these should not be compiled
	    if( child.type && child.type === AT || child.type === AT_COLON ) continue;

	    if(child.vquery){

		this.visitNode(child);

	    } else if(node.mode === MKP){

		this.visitMarkupTok(child, node, i);

	    } else if(node.mode === BLK){

		this.visitBlockTok(child, node, i);

	    } else if(node.mode === EXP){

		this.visitExpressionTok(child, node, i, (nonExp > 0 ? false : true));

	    }
	}

    }

    VCP.getFirstCodeBlock = function(body){
	var i = body.indexOf("BLK({");
	var j = body.indexOf("})BLK");
	if(i != 0 || j == -1) {
	    return ["", body];
	}
	var blk = body.substr(i + 5, j - 5);
	body = body.substr(j + 5);
	return [blk, body];
    }

    VCP.addHead = function(firstCodeBlock, body){
	//todo: should refactor these quick & dirty code
	//      most likely move them to visitNodes
	var params = [];
	var returnType = "string";
	this.sections = [];
	this.layout = "";
	this.imports = {'"bytes"':1};

	//process first code block;
	var lines = firstCodeBlock.split("\n");
	var isImportBlock = false;
	for(var i = 1; i< lines.length; i++) {
	    var l = lines[i].trim();
	    if (l == "") continue;

	    if (l.indexOf('import') == 0) {
		isImportBlock = true
		continue
	    }

	    // End of import
	    if (l == ")") {
		isImportBlock = false
		continue
	    }

	    if (isImportBlock){
		var parts = l.split("/");
		if (parts[parts.length - 2] == "layout") {
		    var layout = parts[parts.length - 1];

		    //Capitalize first character, and ignore '"' at the end
		    this.layout = layout.substr(0, 1).toUpperCase() + layout.substr(1, layout.length - 2);
		    this.imports[l.substr(0, l.length - this.layout.length - 2) + '"'] = 1;
		} else {
		    this.imports[l] = 1;
		}
	    } else if (l.indexOf("var ") == 0 ){
		params.push(l.substring(4));
	    } else {
		console.log("Error Processing: " + this.options["package"] + "/" + this.options["name"] + ".gohtml");
		console.log("Unexpectd: lines in first code block: " + l);
		return;
	    }
	}

	if (body.indexOf("gorazor.HTMLEscape(") > 0) {
	    this.imports['"gorazor"'] = 1;
	}

	imports = Object.keys(this.imports).join("\n");
	params = params.join(", ");

	lines = body.split("\n");

	var inSection = false;
	var counter = 0;
	for(var i=0; i< lines.length; i++) {
	    var l = lines[i].trim();
	    if(l.indexOf("section ") == 0 && l[l.length -1] == "{") {
		sectionName = l.substr(8, l.length - 9).trim();
		this.sections.push(sectionName);
		lines[i] = sectionName + " := func() string {" +
		    "\nvar _buffer bytes.Buffer";
		inSection = true;
		continue;
	    }

	    if(l.indexOf("section ") == 0 && l.indexOf("{") > 0 && l[l.length -1] == "}") {
		sectionName = l.substr(8, l.indexOf("{") - 8).trim();
		this.sections.push(sectionName);
		lines[i] = sectionName + " := func() string {" +
		    "\nreturn ``\n}";
		continue;
	    }

	    if (l == "}" && inSection == true) {
		if (counter == 0) {
		    lines[i] = "return _buffer.String()\n}";
		    inSection = false;
		} else {
		    counter -= 1;
		}
		continue;
	    }

	    if (inSection == true && l[l.length -1] == "{") {
		counter += 1;
	    }
	    if (inSection == true && l[0] == "}") {
		counter -= 1;
	    }
	}

	var newBody = [];
	for(var i=0; i< lines.length; i++) {
	    var l = lines[i].trim();
	    if (l.indexOf("_buffer.WriteString") == 0) {
		l = l.replace("\\'", "'");
	    }
	    if(l != "" & (!l.match(/_buffer.WriteString\("(\\n)+"\)/))) {
		newBody.push(l);
	    }
	}

	body = newBody.join("\n");

	// for (var i = 0; i < this.sections.length; i ++) {
	// 	returnType += ", string";
	// }


	var head = 'package ' + this.options["package"] + '\n\
\n\
import (\n' +
	    imports +'\n)\n\
\n\
func ' + this.options["name"] + '(' + params + ') (' + returnType + ') {\n\
var _buffer bytes.Buffer\n';



	return head + body;
    }

    VCP.addFoot = function(body){
	var foot = '\nreturn ';
	if(this.layout != "") {
	    foot += "layout." + this.layout + "(";
	}
	foot += '_buffer.String()';

	for (var i = 0; i < this.sections.length; i ++) {
	    foot += ", " + this.sections[i] + "()";
	}

	if(this.layout != "") {
	    foot += ")";
	}

	foot += '\n}\n';

	return body + foot;
    }

    VCP.generate = function(){
	var options = this.options;

	// clear whatever's in the current buffer
	this.buffer.length = 0;

	//console.log("ast: ", this.ast)
	this.visitNode(this.ast);

	// coalesce markup
	var joined = this.buffer
	    .join("")
	    .split(")BLKBLK(").join('')
	    .split(")MKPMKP(").join('')
	    .split("MKP(").join( '\n_buffer.WriteString("')
	    .split(")MKP").join('")\n');

	var data = this.getFirstCodeBlock(joined);
	var firstCodeBlock = data[0];
	var body = data[1];

	//support @switch ...{ syntax
	var i = body.indexOf("BLK({");
	while(i > -1) {
	    body = body.substr(0, i) + body.substr(i + 5);
	    i = body.indexOf("})BLK", i);
	    body = body.substr(0, i) + body.substr(i + 5);
	    i = body.indexOf("BLK({");
	}
	body = body.split("BLK(").join('')
	    .split(")BLK").join('');

	joined = this.addHead( firstCodeBlock, body );
	joined = this.addFoot( joined );

	return joined;
    }

    VCompiler.findNonExp = function(node){

	if(node.vquery && node.mode === EXP){
	    return node.filter(VCompiler.findNonExp).length > 0;
	}

	if(node.vquery && node.mode !== EXP){
	    return true;
	} else {
	    return false;
	}
    }

    exports["compile"] = function compile(markup, options){
	var  l
	,tok
	,tokens = []
	,p
	,c
	,cmp
	,i;

	l = new VLexer(markup);
	while(tok = l.advance()) { tokens.push(tok); }

        //console.log("tokens:", tokens)
	tokens.reverse(); // parser needs in reverse order for faster popping vs shift

	p = new VParser(tokens, options);
	p.parse();

	c = new VCompiler(p.ast, markup, options);

	console.log(p.ast)

	cmp = c.generate();
	return cmp;
    };

    return exports;
}({ "version": "0.7.11-3" }));
