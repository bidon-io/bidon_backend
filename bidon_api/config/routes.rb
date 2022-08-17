Rails.application.routes.draw do
  post 'config', to: 'config#create'
end
