Rails.application.routes.draw do
  post 'config', to: 'config#create'
  post 'auction', to: 'auction#create'
end
