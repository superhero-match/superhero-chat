/*
  Copyright (C) 2019 - 2022 MWSOFT
  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.
  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.
  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package config

// App holds the configuration values for the application.
type App struct {
	Port       string `env:"SUPERHERO_CHAT_APP_PORT" yaml:"port" default:":5000"`
	TimeFormat string `env:"SUPERHERO_CHAT_APP_TIME_FORMAT" yaml:"time_format" default:"2006-01-02T15:04:05"`
}
