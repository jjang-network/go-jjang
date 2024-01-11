Pod::Spec.new do |spec|
  spec.name         = 'Gjang'
  spec.version      = '{{.Version}}'
  spec.license      = { :type => 'GNU Lesser General Public License, Version 3.0' }
  spec.homepage     = 'https://github.com/jjang-network/go-jjang'
  spec.authors      = { {{range .Contributors}}
		'{{.Name}}' => '{{.Email}}',{{end}}
	}
  spec.summary      = 'iOS Jjang Client'
  spec.source       = { :git => 'https://github.com/jjang-network/go-jjang.git', :commit => '{{.Commit}}' }

	spec.platform = :ios
  spec.ios.deployment_target  = '9.0'
	spec.ios.vendored_frameworks = 'Frameworks/Gjang.framework'

	spec.prepare_command = <<-CMD
    curl https://gjangstore.blob.core.windows.net/builds/{{.Archive}}.tar.gz | tar -xvz
    mkdir Frameworks
    mv {{.Archive}}/Gjang.framework Frameworks
    rm -rf {{.Archive}}
  CMD
end
