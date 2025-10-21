if &cp | set nocp | endif
let s:cpo_save=&cpo
set cpo&vim
xmap gx <Plug>(open-word-under-cursor)
nmap gx <Plug>(open-word-under-cursor)
xnoremap <Plug>(open-word-under-cursor) <ScriptCmd>vim9.Open(getregion(getpos('v'), getpos('.'), { type: mode() })->join())
nnoremap <Plug>(open-word-under-cursor) <ScriptCmd>vim9.Open(GetWordUnderCursor())
let &cpo=s:cpo_save
unlet s:cpo_save
set autoindent
set background=dark
set backupdir=~/.cache/vim/backup//
set directory=~/.cache/vim/swap//
set fileencodings=ucs-bom,utf-8,default,latin1
set helplang=en
set hlsearch
set ruler
set runtimepath=~/.vim,~/.vim/pack/plugins/start/vim-go,/usr/share/vim/vimfiles,/usr/share/vim/vim91,/usr/share/vim/vim91/pack/dist/opt/netrw,/usr/share/vim/vimfiles/after,~/.vim/after
set shiftwidth=2
set suffixes=.bak,~,.o,.info,.swp,.aux,.bbl,.blg,.brf,.cb,.dvi,.idx,.ilg,.ind,.inx,.jpg,.log,.out,.png,.toc
set tabstop=2
set undodir=~/.cache/vim/undo//
" vim: set ft=vim :
