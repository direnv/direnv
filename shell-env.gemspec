Gem::Specification.new do |s|
  s.name = 'shell-env'
  s.version = '0.' + `git log --oneline | wc -l`.strip
  s.date = Time.at `git log -1 --date=raw | grep Date: | head -n 1 | sed 's/Date: *//'`.to_i
  s.summary = 'path-based environment loading'
  s.homepage = 'http://github.com/zimbatm/shell-env'
  s.author = 'Jonas Pfenniger'
  s.email = 'jonas@pfenniger.name'

  s.files = Dir['bin/**/*'] + Dir['lib/**/*.rb'] + %w( README.md )
  s.executables = %w( shell-env )

  s.has_rdoc = false

  s.add_development_dependency 'ronn'
  s.add_development_dependency 'rake'
end
