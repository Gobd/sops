package age

import (
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMasterKeysFromRecipientsEmpty(t *testing.T) {
	assert := assert.New(t)

	commaSeparatedRecipients := ""
	recipients, err := MasterKeysFromRecipients(commaSeparatedRecipients)

	assert.NoError(err)

	assert.Equal(recipients, make([]*MasterKey, 0))
}

func TestMasterKeyFromRecipientWithLeadingAndTrailingSpaces(t *testing.T) {
	assert := assert.New(t)

	key, err := MasterKeyFromRecipient("  age1yt3tfqlfrwdwx0z0ynwplcr6qxcxfaqycuprpmy89nr83ltx74tqdpszlw  ")

	assert.NoError(err)

	assert.Equal(key.Recipient, "age1yt3tfqlfrwdwx0z0ynwplcr6qxcxfaqycuprpmy89nr83ltx74tqdpszlw")
}

func TestMasterKeysFromRecipientsWithMultiple(t *testing.T) {
	assert := assert.New(t)

	keys, err := MasterKeysFromRecipients("age1yt3tfqlfrwdwx0z0ynwplcr6qxcxfaqycuprpmy89nr83ltx74tqdpszlw,age1tmaae3ld5vpevmsh5yacsauzx8jetg300mpvc4ugp5zr5l6ssq9sla97ep")

	assert.NoError(err)

	assert.Equal(len(keys), 2)
	assert.Equal(keys[0].Recipient, "age1yt3tfqlfrwdwx0z0ynwplcr6qxcxfaqycuprpmy89nr83ltx74tqdpszlw")
	assert.Equal(keys[1].Recipient, "age1tmaae3ld5vpevmsh5yacsauzx8jetg300mpvc4ugp5zr5l6ssq9sla97ep")
}

func TestAge(t *testing.T) {
	assert := assert.New(t)

	keys, err := MasterKeysFromRecipients("age1yt3tfqlfrwdwx0z0ynwplcr6qxcxfaqycuprpmy89nr83ltx74tqdpszlw,age1tmaae3ld5vpevmsh5yacsauzx8jetg300mpvc4ugp5zr5l6ssq9sla97ep")

	assert.NoError(err)
	assert.Equal(len(keys), 2)
	assert.Equal(keys[0].Recipient, "age1yt3tfqlfrwdwx0z0ynwplcr6qxcxfaqycuprpmy89nr83ltx74tqdpszlw")
	assert.Equal(keys[1].Recipient, "age1tmaae3ld5vpevmsh5yacsauzx8jetg300mpvc4ugp5zr5l6ssq9sla97ep")

	dataKey := []byte("abcdefghijklmnopqrstuvwxyz123456")

	for _, key := range keys {
		err = key.Encrypt(dataKey)
		assert.NoError(err)

		_, filename, _, _ := runtime.Caller(0)
		err = os.Setenv("SOPS_AGE_KEY_FILE", path.Join(path.Dir(filename), "keys.txt"))
		assert.NoError(err)

		decryptedKey, err := key.Decrypt()
		assert.NoError(err)
		assert.Equal(dataKey, decryptedKey)
	}

}

func TestAgeDotEnv(t *testing.T) {
	assert := assert.New(t)

	key, err := MasterKeyFromRecipient("age1yt3tfqlfrwdwx0z0ynwplcr6qxcxfaqycuprpmy89nr83ltx74tqdpszlw")

	assert.NoError(err)
	assert.Equal("age1yt3tfqlfrwdwx0z0ynwplcr6qxcxfaqycuprpmy89nr83ltx74tqdpszlw", key.ToString())

	dotenv := `IMAGE_PREFIX=repo/service-
APPLICATION_KEY=K6pfAWuUVND9Fz5SC7jmA6pfAWuUVND9Fz5SC7jmA
KEY_ID=003683d721f2ae683d721f2a1
DOMAIN=files.127.0.0.1.nip.io`
	dataKey := []byte(dotenv)

	err = key.Encrypt(dataKey)
	assert.NoError(err)

	_, filename, _, _ := runtime.Caller(0)
	err = os.Setenv(SopsAgeKeyFileEnv, path.Join(path.Dir(filename), "keys.txt"))
	defer os.Unsetenv(SopsAgeKeyFileEnv)
	assert.NoError(err)

	decryptedKey, err := key.Decrypt()
	assert.NoError(err)
	assert.Equal(dataKey, decryptedKey)
}

func TestAgeEnv(t *testing.T) {
	assert := assert.New(t)

	key, err := MasterKeyFromRecipient("age1yt3tfqlfrwdwx0z0ynwplcr6qxcxfaqycuprpmy89nr83ltx74tqdpszlw")

	assert.NoError(err)
	assert.Equal("age1yt3tfqlfrwdwx0z0ynwplcr6qxcxfaqycuprpmy89nr83ltx74tqdpszlw", key.ToString())

	dataKey := []byte("abcdefghijklmnopqrstuvwxyz123456")

	err = key.Encrypt(dataKey)
	assert.NoError(err)

	_, filename, _, _ := runtime.Caller(0)
	keysBytes, err := ioutil.ReadFile(path.Join(path.Dir(filename), "keys.txt"))
	assert.NoError(err)
	err = os.Setenv(SopsAgeKeyEnv, string(keysBytes))
	defer os.Unsetenv(SopsAgeKeyEnv)
	assert.NoError(err)

	decryptedKey, err := key.Decrypt()
	assert.NoError(err)
	assert.Equal(dataKey, decryptedKey)
}
