#!/usr/bin/env ruby
load File.expand_path('../bin/shell-env', __FILE__)
require "test/unit"

class TestMe < Test::Unit::TestCase
  include ShellEnv
    
  def test_diff_env
    assert_equal({'one'=>nil}, diff_env({'one'=>'x'}, {}))
    assert_equal({'two'=>'2'}, diff_env({}, {'two'=>'2'}))
  end
    
  def test_marshal
    env = {'one'=>'1', 'two'=>'2'}
    assert_equal(env, unmarshal_env(marshal_env(env)))
  end
end
