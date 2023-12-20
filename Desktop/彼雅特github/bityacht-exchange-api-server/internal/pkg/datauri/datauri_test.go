package datauri

import (
	"bityacht-exchange-api-server/configs"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/logger"
	"testing"

	"github.com/spf13/viper"
)

func TestValidateImage(t *testing.T) {
	tests := []struct {
		dataURI     string
		wantErrCode errpkg.Code
	}{
		{`data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABgAAAAYCAYAAADgdz34AAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAApgAAAKYB3X3/OAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAANCSURBVEiJtZZPbBtFFMZ/M7ubXdtdb1xSFyeilBapySVU8h8OoFaooFSqiihIVIpQBKci6KEg9Q6H9kovIHoCIVQJJCKE1ENFjnAgcaSGC6rEnxBwA04Tx43t2FnvDAfjkNibxgHxnWb2e/u992bee7tCa00YFsffekFY+nUzFtjW0LrvjRXrCDIAaPLlW0nHL0SsZtVoaF98mLrx3pdhOqLtYPHChahZcYYO7KvPFxvRl5XPp1sN3adWiD1ZAqD6XYK1b/dvE5IWryTt2udLFedwc1+9kLp+vbbpoDh+6TklxBeAi9TL0taeWpdmZzQDry0AcO+jQ12RyohqqoYoo8RDwJrU+qXkjWtfi8Xxt58BdQuwQs9qC/afLwCw8tnQbqYAPsgxE1S6F3EAIXux2oQFKm0ihMsOF71dHYx+f3NND68ghCu1YIoePPQN1pGRABkJ6Bus96CutRZMydTl+TvuiRW1m3n0eDl0vRPcEysqdXn+jsQPsrHMquGeXEaY4Yk4wxWcY5V/9scqOMOVUFthatyTy8QyqwZ+kDURKoMWxNKr2EeqVKcTNOajqKoBgOE28U4tdQl5p5bwCw7BWquaZSzAPlwjlithJtp3pTImSqQRrb2Z8PHGigD4RZuNX6JYj6wj7O4TFLbCO/Mn/m8R+h6rYSUb3ekokRY6f/YukArN979jcW+V/S8g0eT/N3VN3kTqWbQ428m9/8k0P/1aIhF36PccEl6EhOcAUCrXKZXXWS3XKd2vc/TRBG9O5ELC17MmWubD2nKhUKZa26Ba2+D3P+4/MNCFwg59oWVeYhkzgN/JDR8deKBoD7Y+ljEjGZ0sosXVTvbc6RHirr2reNy1OXd6pJsQ+gqjk8VWFYmHrwBzW/n+uMPFiRwHB2I7ih8ciHFxIkd/3Omk5tCDV1t+2nNu5sxxpDFNx+huNhVT3/zMDz8usXC3ddaHBj1GHj/As08fwTS7Kt1HBTmyN29vdwAw+/wbwLVOJ3uAD1wi/dUH7Qei66PfyuRj4Ik9is+hglfbkbfR3cnZm7chlUWLdwmprtCohX4HUtlOcQjLYCu+fzGJH2QRKvP3UNz8bWk1qMxjGTOMThZ3kvgLI5AzFfo379UAAAAASUVORK5CYII=`, 0},
		{`data:image/gif;base64,R0lGODlhAQABAAAAACw=`, errpkg.CodeBadUploadedFileType},
		{`data:image/jpeg;base64,R0lGODlhAQABAAAAACw=`, errpkg.CodeBadImageData},
		{`data:image/png;base64,R0lGODlhAQABAAAAACw=`, errpkg.CodeBadImageData},
		{`data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAMCAgICAgMCAgIDAwMDBAYEBAQEBAgGBgUGCQgKCgkICQkKDA8MCgsOCwkJDRENDg8QEBEQCgwSExIQEw8QEBD/2wBDAQMDAwQDBAgEBAgQCwkLEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBD/wAARCAAUAGQDASIAAhEBAxEB/8QAFwABAQEBAAAAAAAAAAAAAAAAAAUHBv/EABoQAQACAwEAAAAAAAAAAAAAAAAGRIKDwcL/xAAbAQACAgMBAAAAAAAAAAAAAAAGBwAIAQQFCf/EACERAAEBCAMBAAAAAAAAAAAAAAAGAQIDBTWBgrI0Q8Gx/9oADAMBAAIRAxEAPwDcwDIAc6iE3dfp1Dl4Td1+nUEqrazGx0dHIlKRBy2eDVWVNVbKX7cfTM967+BGkVfPiyjSKvnx2ptw37fWAMoabEtswigAsWZjoCt5Z0AIQngPQsQR1EJu6/TqAJVW1mNjo6ORKUiDls8Gqg2Uv24+mZ7138CNIq+fAdqbcN+31gDKGmxLbMIoALFmY6AreWdACEP/2Q==`, 0},
		{`data:image/png;base64,`, errpkg.CodeBadImageData},
		{`data:image/jpeg;base64`, errpkg.CodeBadDataURIImageFormat},
	}

	for _, test := range tests {
		err := ValidateImage(test.dataURI)
		if test.wantErrCode == 0 {
			if err != nil {
				t.Errorf("Wrong ValidateImage Result, Want err == nil but Return: %d(%+v)", err.Code, err.Err)
			}
		} else if err == nil {
			t.Errorf("Wrong ValidateImage Result, Want err.Code == %d but Return nil", test.wantErrCode)
		} else if err.Code != test.wantErrCode {
			t.Errorf("Wrong ValidateImage Result, Want err.Code == %d but Return: %d(%+v)", test.wantErrCode, err.Code, err.Err)
		}
	}
}

func TestDownloadImage(t *testing.T) {
	viper.AddConfigPath("../../../configs")
	configs.Init()
	logger.Init()

	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "jpg", args: args{url: `https://www.whitestone-gallery.com/cdn/shop/files/20220813_tw_yyl_top-sp.jpg?v=1660547673`}, want: 357563},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DownloadImage(logger.Logger, tt.args.url); len(got) != tt.want {
				t.Errorf("len(DownloadImage()) = %v, want %v", len(got), tt.want)
			}
		})
	}
}
