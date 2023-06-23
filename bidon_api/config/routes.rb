Rails.application.routes.draw do
  get 'health_checks', to: 'health#show'

  post 'config', to: 'config#create'

  post 'auction/:ad_type', to: 'auction#create'
  post 'bidding/:ad_type', to: 'bidding#create'
  post 'stats/:ad_type',   to: 'stats#create'
  post 'click/:ad_type',   to: 'click#create'
  post 'show/:ad_type',    to: 'show#create'
  post 'loss/:ad_type',    to: 'loss#create'
  post 'win/:ad_type',     to: 'win#create'
  post 'reward/rewarded',  to: 'reward#create'

  post ':ad_type/auction', to: 'auction#create'
  post ':ad_type/stats',   to: 'stats#create'
  post ':ad_type/click',   to: 'click#create'
  post ':ad_type/show',    to: 'show#create'
  post 'rewarded/reward',  to: 'reward#create'
end
