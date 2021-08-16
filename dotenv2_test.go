package dotenv_test

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/fairyhunter13/dotenv"
	"github.com/stretchr/testify/assert"
)

func loadEnvAndCompareValues2(t *testing.T, loader func(opts ...dotenv.FnOption) error, opts []dotenv.FnOption, envFileName string, expectedValues map[string]string, presets map[string]string) {
	// first up, clear the env
	os.Clearenv()

	for k, v := range presets {
		os.Setenv(k, v)
	}

	err := loader(opts...)
	if err != nil {
		t.Fatalf("Error loading %v", envFileName)
	}

	for k := range expectedValues {
		envValue := os.Getenv(k)
		v := expectedValues[k]
		if envValue != v {
			t.Errorf("Mismatch for key '%v': expected '%v' got '%v'", k, v, envValue)
			t.Errorf("rune expected: %b got %b", []byte(v), []byte(envValue))
		}
	}
}

func TestLoadWithNoArgsLoadsDotEnv2(t *testing.T) {
	err := dotenv.Load2()
	pathError := err.(*os.PathError)
	if pathError == nil || pathError.Op != "open" || pathError.Path != ".env" {
		t.Errorf("Didn't try and open .env by default")
	}
}

func TestOverloadWithNoArgsOverloadsDotEnv2(t *testing.T) {
	err := dotenv.Load2(
		dotenv.WithOverload(true),
	)
	pathError := err.(*os.PathError)
	if pathError == nil || pathError.Op != "open" || pathError.Path != ".env" {
		t.Errorf("Didn't try and open .env by default")
	}
}

func TestLoadFileNotFound2(t *testing.T) {
	err := dotenv.Load2(
		dotenv.WithPaths("somefilethatwillneverexistever.env"),
	)
	if err == nil {
		t.Error("File wasn't found but Load didn't return an error")
	}
}

func TestOverloadFileNotFound2(t *testing.T) {
	err := dotenv.Load2(
		dotenv.WithOverload(true),
		dotenv.WithPaths("somefilethatwillneverexistever.env"),
	)
	if err == nil {
		t.Error("File wasn't found but Overload didn't return an error")
	}
}

func TestReadPlainEnv2(t *testing.T) {
	envFileName := "fixtures/plain.env"
	expectedValues := map[string]string{
		"OPTION_A": "1",
		"OPTION_B": "2",
		"OPTION_C": "3",
		"OPTION_D": "4",
		"OPTION_E": "5",
		"OPTION_F": "",
		"OPTION_G": "",
	}

	envMap, err := dotenv.ReadFile2(envFileName)
	if err != nil {
		t.Error("Error reading file")
	}

	if len(envMap) != len(expectedValues) {
		t.Error("Didn't get the right size map back")
	}

	for key, value := range expectedValues {
		if envMap[key] != value {
			t.Error("Read got one of the keys wrong")
		}
	}
}

func TestLoadDoesNotOverride2(t *testing.T) {
	envFileName := "fixtures/plain.env"

	// ensure NO overload
	presets := map[string]string{
		"OPTION_A": "do_not_override",
		"OPTION_B": "",
	}

	expectedValues := map[string]string{
		"OPTION_A": "do_not_override",
		"OPTION_B": "",
	}
	loadEnvAndCompareValues2(
		t,
		dotenv.Load2,
		[]dotenv.FnOption{
			dotenv.WithPaths(envFileName),
		},
		envFileName,
		expectedValues,
		presets,
	)
}

func TestOverloadDoesOverride2(t *testing.T) {
	envFileName := "fixtures/plain.env"

	// ensure NO overload
	presets := map[string]string{
		"OPTION_A": "do_not_override",
	}

	expectedValues := map[string]string{
		"OPTION_A": "1",
	}
	loadEnvAndCompareValues2(
		t,
		dotenv.Load2,
		[]dotenv.FnOption{
			dotenv.WithOverload(true),
			dotenv.WithPaths(envFileName),
		},
		envFileName,
		expectedValues,
		presets,
	)
}

func TestLoadPlainEnv2(t *testing.T) {
	envFileName := "fixtures/plain.env"
	expectedValues := map[string]string{
		"OPTION_A": "1",
		"OPTION_B": "2",
		"OPTION_C": "3",
		"OPTION_D": "4",
		"OPTION_E": "5",
	}

	loadEnvAndCompareValues2(
		t,
		dotenv.Load2,
		[]dotenv.FnOption{
			dotenv.WithPaths(envFileName),
		},
		envFileName,
		expectedValues,
		noopPresets,
	)
}

func TestLoadExportedEnv2(t *testing.T) {
	envFileName := "fixtures/exported.env"
	expectedValues := map[string]string{
		"OPTION_A": "2",
		"OPTION_B": "\n",
	}

	loadEnvAndCompareValues2(
		t,
		dotenv.Load2,
		[]dotenv.FnOption{
			dotenv.WithPaths(envFileName),
		},
		envFileName,
		expectedValues,
		noopPresets,
	)
}

func TestLoadEqualsEnv2(t *testing.T) {
	envFileName := "fixtures/equals.env"
	expectedValues := map[string]string{
		"OPTION_A": "postgres://localhost:5432/database?sslmode=disable",
	}

	loadEnvAndCompareValues2(
		t,
		dotenv.Load2,
		[]dotenv.FnOption{
			dotenv.WithPaths(envFileName),
		},
		envFileName,
		expectedValues,
		noopPresets,
	)
}

func TestLoadQuotedEnv2(t *testing.T) {
	envFileName := "fixtures/quoted.env"
	expectedValues := map[string]string{
		"OPTION_A": "1",
		"OPTION_B": "2",
		"OPTION_C": "",
		"OPTION_D": "\n",
		"OPTION_E": "1",
		"OPTION_F": "2",
		"OPTION_G": "",
		"OPTION_H": "\n",
	}

	loadEnvAndCompareValues2(
		t,
		dotenv.Load2,
		[]dotenv.FnOption{
			dotenv.WithPaths(envFileName),
		},
		envFileName,
		expectedValues,
		noopPresets,
	)
}
func TestLoadBenchEnv2(t *testing.T) {
	envFileName := "fixtures/bench.env"
	expectedValues := map[string]string{
		"A": "1",
		"B": "2",
		"C": "3",
		"D": "4",
		"E": "5",
		"F": "SOMETHING",
		"G": "something 'else'",
		"H": "SOMETHING else #2",
		"I": "something escaped\"",
		"J": "asdfa",
		"K": "http",
		"L": "http://",
		"M": "http://github.com",
		"N": "http://github.com/fairyhunter13",
		"O": "http://github.com/fairyhunter13/dotenv",
		"P": "124215",
		"Q": "127.0.0.1",
		"R": ";aklsdgj",
		"S": "adsg;hkjl",
		"T": "k;lajdsg",
		"U": "\n",
		"V": "\n",
		"X": "\r",
		"Y": "\r\n",
		"Z": "\"",
	}

	loadEnvAndCompareValues2(
		t,
		dotenv.Load2,
		[]dotenv.FnOption{
			dotenv.WithPaths(envFileName),
		},
		envFileName,
		expectedValues,
		noopPresets,
	)
}

func TestActualEnvVarsAreLeftAlone2(t *testing.T) {
	os.Clearenv()
	os.Setenv("OPTION_A", "actualenv")
	_ = dotenv.Load2(
		dotenv.WithPaths("fixtures/plain.env"),
	)

	if os.Getenv("OPTION_A") != "actualenv" {
		t.Error("An ENV var set earlier was overwritten")
	}
}

func TestErrorReadDirectory2(t *testing.T) {
	envFileName := "fixtures/"
	envMap, err := dotenv.ReadFile2(envFileName)

	if err == nil {
		t.Errorf("Expected error, got %v", envMap)
	}
}

func TestErrorParsing2(t *testing.T) {
	envFileName := "fixtures/invalid1.env"
	envMap, err := dotenv.ReadFile2(envFileName)
	if err == nil {
		t.Errorf("Expected error, got %v", envMap)
	}
}

func BenchmarkDotenv2(b *testing.B) {
	b.StopTimer()
	f, _ := os.Open("fixtures/bench.env")
	defer f.Close()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = dotenv.Read2(f)
	}
	b.StopTimer()
}

func TestMultilineParsing(t *testing.T) {
	envFileName := "fixtures/multiline.env"
	file, err := os.Open(envFileName)
	assert.Nil(t, err)
	defer file.Close()

	err = dotenv.LoadReader2(file)
	assert.Nil(t, err)
	_, _ = file.Seek(0, io.SeekStart)

	envMap, err := dotenv.Read2(file)
	assert.Nil(t, err)
	assert.NotNil(t, envMap)

	assert.Equal(t, "hello fairy!", envMap["TEST"])

	jsonFlag := map[string]bool{}
	err = json.Unmarshal([]byte(envMap["JSON_FLAG"]), &jsonFlag)
	assert.Nil(t, err)

	assert.Equal(t, true, jsonFlag["value1"])
	assert.Equal(t, false, jsonFlag["value2"])
	assert.Equal(t, true, jsonFlag["value3"])

	jsonString := map[string]string{}
	err = json.Unmarshal([]byte(envMap["JSON_STRING"]), &jsonString)
	assert.Nil(t, err)

	assert.Equal(t, "hello ", jsonString["start"])
	assert.Equal(t, "guys!", jsonString["end"])
}
