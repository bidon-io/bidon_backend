namespace :appodeal do
  task sync_test_apps: :environment do
    Appodeal::SyncData.new.call
  end
end
