# Robo Hash

<p>
  <a href="https://choosealicense.com/licenses/mit/">
    <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="MIT License">
  </a>
  <img src="https://img.shields.io/badge/Go-%20blue?logo=go" alt="Go"  style="margin-left: 10px;">
</p>

This project is a port of [Robohash](https://github.com/e1ven/Robohash) to Go. Follow the steps below to install, build, and package the project.

<br>

## Installation Guide

Follow the steps below to install, build, and package the project.

### Prerequisites

Ensure you have the following installed:

- [Go](https://golang.org/dl/) (Version 1.23 or above)
- [Make](https://www.gnu.org/software/make/)
- [tar](https://www.gnu.org/software/tar/) (for creating tar archives)

### Installation

1. Clone the Repository
2. Build the Project

```bash
make all
cd build
tar -xvzf robohash.tar.gz
```

<br>
<br>

## Usage


### Construction Settings

| Parameter    | Type     | Description                                                                                    |
|:-------------|:---------|:-----------------------------------------------------------------------------------------------|
| `-input`     | `string` | **Required**. Input to be hashed                                                               |
| `-set`       | `string` | - 1 -> Robot<br/>- 2 -> Monster<br/>-3 -> Robot2<br/>4 -> Cat<br/>5 -> Person<br/>- any        |
| `-color`     | `string` | Only for **set1**<br/>[[blue, brown, green, grey,  orange, pink, purple, red, white, yellow]   |
| `-extension` | `string` | File to be saved (datauri will not be saved to file)<br/>  [png, jpg, jpeg, gif, ppm, datauri] |
| `-bgset`     | `string` | Adds a background<br/> [1,2,any]                                                               |
| `-sizex`     | `int`    | X size in px                                                                                   |
| `-sizey`     | `int`    | Y size in px                                                                                   |

<br>

### Hash Setting
| Parameter    | Type   | Description                                                                      |
|:-------------|:-------|:---------------------------------------------------------------------------------|
| `-ignoreExt` | `none` | It will ignore image extensions, such as .png, .jpg, etc. When creating the hash |
| `-slots`     | `int`  | Number of slots the hash will be divided. **Be careful with this option**        |

<br>

> `any` -> refers to random value

<br>
<br>

## Examples

<img src="docs/robo1.png">
<img src="docs/robo2.png">
<img src="docs/robo3.png">

<br>
<br>

## Attribution

This project uses images from [Robohash.org](https://robohash.org). The images are created by:

- **Set 1**: Robots by Zikri Kader
- **Set 2**: Robots by Hrvoje Novakovic
- **Set 3**: Robots by Julian Peter Arias
- **Set 4**: Cats by David Revoy
- **Set 5**: Avatars by Pablo Stanley

You are free to embed these images under the terms of the [CC-BY License](https://creativecommons.org/licenses/by/4.0/).
