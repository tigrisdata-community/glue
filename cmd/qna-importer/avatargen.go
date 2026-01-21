package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"path"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/tigrisdata-community/glue/web/sdcpp"
	"github.com/tigrisdata/storage-go"
)

// SHA256sum computes a cryptographic hash. Still used for proof-of-work challenges
// where we need the security properties of a cryptographic hash function.
func SHA256sum(text string) string {
	hash := sha256.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

type AvatarGen struct {
	sd     *sdcpp.Client
	tigris *storage.Client
	bucket string
	prefix string
}

func (a *AvatarGen) GenerateAndUpload(ctx context.Context, input string) (string, error) {
	hash := SHA256sum(input)
	prompt, _ := a.hallucinatePrompt(hash)

	resp, err := a.sd.Generate(ctx, sdcpp.ImageGenerationRequest{
		Prompt:       prompt,
		Size:         "512x512",
		OutputFormat: "png",
	})
	if err != nil {
		return "", fmt.Errorf("can't generate image: %w", err)
	}

	data, err := resp.ToWebP(0, 75)
	if err != nil {
		return "", fmt.Errorf("can't encode image: %w", err)
	}

	key := path.Join(a.prefix, hash+".webp")

	if _, err := a.tigris.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(a.bucket),
		Key:           aws.String(key),
		Body:          io.NopCloser(bytes.NewBuffer(data)),
		ACL:           types.ObjectCannedACLPublicRead,
		ContentType:   aws.String("image/webp"),
		ContentLength: aws.Int64(int64(len(data))),
	}); err != nil {
		return "", fmt.Errorf("can't upload object: %w", err)
	}

	return key, nil
}

func (a *AvatarGen) hallucinatePrompt(hash string) (string, int) {
	var sb strings.Builder
	fmt.Fprint(&sb, "headshot, ")
	if hash[0] > '0' && hash[0] <= '5' {
		fmt.Fprint(&sb, "1girl, ")
	} else {
		fmt.Fprint(&sb, "1guy, ")
	}

	switch hash[1] {
	case '0', '1', '2', '3':
		fmt.Fprint(&sb, "blonde, ")
	case '4', '5', '6', '7':
		fmt.Fprint(&sb, "brown hair, ")
	case '8', '9', 'a', 'b':
		fmt.Fprint(&sb, "red hair, ")
	case 'c', 'd', 'e', 'f':
		fmt.Fprint(&sb, "black hair, ")
	default:
	}

	if hash[2] > '0' && hash[2] <= '5' {
		fmt.Fprint(&sb, "coffee shop, ")
	} else {
		fmt.Fprint(&sb, "landscape, outdoors, ")
	}

	if hash[3] > '0' && hash[3] <= '5' {
		fmt.Fprint(&sb, "hoodie, ")
	} else {
		fmt.Fprint(&sb, "sweatsuit, ")
	}

	switch hash[4] {
	case '0', '1', '2', '3':
		fmt.Fprint(&sb, "final fantasy 14, ")
	case '4', '5', '6', '7':
		fmt.Fprint(&sb, "breath of the wild, ")
	case '8', '9', 'a', 'b':
		fmt.Fprint(&sb, "genshin impact, ")
	case 'c', 'd', 'e', 'f':
		fmt.Fprint(&sb, "arknights, ")
	default:
	}

	if hash[5] > '0' && hash[5] <= '5' {
		fmt.Fprint(&sb, "watercolor, ")
	} else {
		fmt.Fprint(&sb, "matte painting, ")
	}

	switch hash[6] {
	case '0', '1', '2', '3':
		fmt.Fprint(&sb, "highly detailed, ")
	case '4', '5', '6', '7':
		fmt.Fprint(&sb, "ornate, ")
	case '8', '9', 'a', 'b':
		fmt.Fprint(&sb, "thick lines, ")
	case 'c', 'd', 'e', 'f':
		fmt.Fprint(&sb, "3d render, ")
	default:
	}

	switch hash[7] {
	case '0', '1', '2', '3':
		fmt.Fprint(&sb, "short hair, ")
	case '4', '5', '6', '7':
		fmt.Fprint(&sb, "long hair, ")
	case '8', '9', 'a', 'b':
		fmt.Fprint(&sb, "ponytail, ")
	case 'c', 'd', 'e', 'f':
		fmt.Fprint(&sb, "pigtails, ")
	default:
	}

	switch hash[8] {
	case '0', '1', '2', '3':
		fmt.Fprint(&sb, "smile, ")
	case '4', '5', '6', '7':
		fmt.Fprint(&sb, "frown, ")
	case '8', '9', 'a', 'b':
		fmt.Fprint(&sb, "laughing, ")
	case 'c', 'd', 'e', 'f':
		fmt.Fprint(&sb, "angry, ")
	default:
	}

	switch hash[9] {
	case '0', '1', '2', '3':
		fmt.Fprint(&sb, "sweater, ")
	case '4', '5', '6', '7':
		fmt.Fprint(&sb, "tshirt, ")
	case '8', '9', 'a', 'b':
		fmt.Fprint(&sb, "suitjacket, ")
	case 'c', 'd', 'e', 'f':
		fmt.Fprint(&sb, "armor, ")
	default:
	}

	switch hash[10] {
	case '0', '1', '2', '3':
		fmt.Fprint(&sb, "blue eyes, ")
	case '4', '5', '6', '7':
		fmt.Fprint(&sb, "red eyes, ")
	case '8', '9', 'a', 'b':
		fmt.Fprint(&sb, "brown eyes, ")
	case 'c', 'd', 'e', 'f':
		fmt.Fprint(&sb, "hazel eyes, ")
	default:
	}

	if hash[11] == '0' {
		fmt.Fprint(&sb, "heterochromia, ")

		switch hash[10] {
		case '0', '1', '2', '3':
			fmt.Fprint(&sb, "red eyes, ")
		case '4', '5', '6', '7':
			fmt.Fprint(&sb, "yellow eyes, ")
		case '8', '9', 'a', 'b':
			fmt.Fprint(&sb, "purple eyes, ")
		case 'c', 'd', 'e', 'f':
			fmt.Fprint(&sb, "green eyes, ")
		default:
		}
	}

	switch hash[12] {
	case '0', '1', '2', '3':
		fmt.Fprint(&sb, "morning, ")
	case '4', '5', '6', '7':
		fmt.Fprint(&sb, "afternoon, ")
	case '8', '9', 'a', 'b':
		fmt.Fprint(&sb, "evening, ")
	case 'c', 'd', 'e', 'f':
		fmt.Fprint(&sb, "nighttime, ")
	default:
	}

	switch hash[14] {
	case '0', '1', '2', '3':
		fmt.Fprint(&sb, "pixar, ")
	case '4', '5', '6', '7':
		fmt.Fprint(&sb, "anime, ")
	case '8', '9', 'a', 'b':
		fmt.Fprint(&sb, "studio ghibli, ")
	case 'c', 'd', 'e', 'f':
		fmt.Fprint(&sb, "LinkedIn, ")
	default:
	}

	seedPortion := hash[len(hash)-9 : len(hash)-1]
	seed, err := strconv.ParseInt(seedPortion, 16, 32)
	if err != nil {
		seed = int64(rand.Int())
	}

	fmt.Fprint(&sb, "pants")

	return sb.String(), int(seed)
}
