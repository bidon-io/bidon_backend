Rails.application.routes.draw do
  post 'config', to: 'config#create'
  post 'auction/:ad_type', to: 'auction#create'

  post 'stats/:ad_type',   to: 'stats#create'
  post 'click/:ad_type',   to: 'click#create'
  post 'show/:ad_type',    to: 'show#create'
  post 'reward/:ad_type',  to: 'reward#create'
end
