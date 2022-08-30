Rails.application.routes.draw do
  post 'config',  to: 'config#create'
  post 'auction', to: 'auction#create'
  post 'stats',   to: 'stats#create'
  post 'click',   to: 'click#create'
  post 'finish',  to: 'finish#create'
  post 'show',    to: 'show#create'
end
