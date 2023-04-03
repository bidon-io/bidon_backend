Rails.application.routes.draw do
  get 'health_checks', to: 'health#show'

  mount Avo::Engine, at: Avo.configuration.root_path
end
