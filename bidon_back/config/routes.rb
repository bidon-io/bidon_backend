Rails.application.routes.draw do
  get 'up', to: 'health#show'

  mount Avo::Engine, at: Avo.configuration.root_path
end
