workers Integer(ENV['PUMA_WORKER_COUNT'] || 2)
threads_count = Integer(ENV['PUMA_THREADS_COUNT'] || 5)
threads threads_count, threads_count

preload_app!

port ENV.fetch('PORT', 3000)
environment ENV.fetch('RAILS_ENV', 'development')
pidfile ENV.fetch('PIDFILE', 'tmp/pids/server.pid')

# Allow puma to be restarted by `bin/rails restart` command.
plugin :tmp_restart
