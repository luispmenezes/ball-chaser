package header

import (
	"errors"
	"time"
)

type Header struct {
	Size            uint
	CRC             uint
	EngineVersion   uint
	LicenseeVersion uint
	NetVersion      uint
	Label           string
	Properties      map[string]Property
}

func (h *Header) GetBuildVersion() (string, error) {
	prop, propExists := h.Properties[buildVersionProperty]

	if !propExists {
		return "", errors.New("build version not found in properties")
	}

	return prop.getStringValue()
}

func (h *Header) GetBuildId() (string, error) {
	prop, propExists := h.Properties[buildIDProperty]

	if !propExists {
		return "", errors.New("build id not found in properties")
	}

	return prop.getStringValue()
}

func (h *Header) GetGameVersion() (string, error) {
	prop, propExists := h.Properties[gameVersionProperty]

	if !propExists {
		return "", errors.New("game version not found in properties")
	}

	return prop.getStringValue()
}

func (h *Header) GetReplayId() (string, error) {
	prop, propExists := h.Properties[replayIdProperty]

	if !propExists {
		return "", errors.New("replay id not found in properties")
	}

	return prop.getStringValue()
}

func (h *Header) GetMapName() (string, error) {
	prop, propExists := h.Properties[mapNameProperty]

	if !propExists {
		return "", errors.New("map name not found in properties")
	}

	return prop.getStringValue()
}

func (h *Header) GetDate() (time.Time, error) {
	prop, propExists := h.Properties[dateProperty]

	if !propExists {
		return time.Time{}, errors.New("date not found in properties")
	}

	strVal, err := prop.getStringValue()

	if err != nil {
		return time.Time{}, err
	}

	date, err := time.Parse(datePropertyFormat, strVal)

	if err != nil {
		return time.Time{}, errors.New("invalid date format")
	}

	return date, nil
}

func (h *Header) GetNumberOfFrames() (uint, error) {
	prop, propExists := h.Properties[numberOfFramesProperty]

	if !propExists {
		return 0, errors.New("number of frames not found in properties")
	}

	return prop.getIntValue()
}

func (h *Header) GetMatchType() (string, error) {
	prop, propExists := h.Properties[matchTypeProperty]

	if !propExists {
		return "", errors.New("match type not found in properties")
	}

	return prop.getStringValue()
}

func (h *Header) GetTeamSize() (uint, error) {
	prop, propExists := h.Properties[teamSizeProperty]

	if !propExists {
		return 0, errors.New("team size not found in properties")
	}

	return prop.getIntValue()
}

func (h *Header) GetUnfairTeamSize() (uint, error) {
	prop, propExists := h.Properties[unfairTeamSizeProperty]

	if !propExists {
		return 0, errors.New("unfair team size not found in properties")
	}

	return prop.getIntValue()
}

func (h *Header) GetTeam1Score() (uint, error) {
	prop, propExists := h.Properties[team1ScoreProperty]

	if !propExists {
		return 0, nil
	}

	return prop.getIntValue()
}

func (h *Header) GetTeam2Score() (uint, error) {
	prop, propExists := h.Properties[team2ScoreProperty]

	if !propExists {
		return 0, nil
	}

	return prop.getIntValue()
}

func (h *Header) GetGoals() ([]Goal, error) {
	var goals []Goal

	prop, propExists := h.Properties[goalsProperty]

	if !propExists {
		return goals, nil
	}

	arrayProps := prop.Value.([]map[string]Property)

	for _, arrP := range arrayProps {
		frameProp := arrP[goalsPropertyFrameKey]
		frame, err := frameProp.getIntValue()
		if err != nil {
			return nil, err
		}
		playerProp := arrP[goalsPropertyPlayerKey]
		player, err := playerProp.getStringValue()
		if err != nil {
			return nil, err
		}
		teamProp := arrP[goalsPropertyTeamKey]
		team, err := teamProp.getIntValue()
		if err != nil {
			return nil, err
		}
		goals = append(goals, Goal{
			Frame:  frame,
			Scorer: player,
			Team:   team,
		})
	}

	return goals, nil
}

func (h *Header) GetHighlights() ([]Highlight, error) {
	var highlights []Highlight

	prop, propExists := h.Properties[highlightsProperty]

	if !propExists {
		return highlights, nil
	}

	arrayProps := prop.Value.([]map[string]Property)

	for _, arrP := range arrayProps {
		frameProp := arrP[highlightsPropertyFrameKey]
		frame, err := frameProp.getIntValue()
		if err != nil {
			return nil, err
		}
		carProp := arrP[highlightsPropertyCarKey]
		car, err := carProp.getStringValue()
		if err != nil {
			return nil, err
		}
		ballProp := arrP[highlightsPropertyBallKey]
		ball, err := ballProp.getStringValue()
		if err != nil {
			return nil, err
		}
		highlights = append(highlights, Highlight{
			Frame: frame,
			Car:   car,
			Ball:  ball,
		})
	}

	return highlights, nil
}

func (h *Header) GetPlayerStatistics() ([]PlayerStat, error) {
	var stats []PlayerStat

	prop, propExists := h.Properties[playerStatsProperty]

	if !propExists {
		return stats, nil
	}

	arrayProps := prop.Value.([]map[string]Property)

	for _, arrP := range arrayProps {
		nameProp := arrP[playerStatsPropertyNameKey]
		name, err := nameProp.getStringValue()
		if err != nil {
			return nil, err
		}
		plaformKeyProp := arrP[playerStatsPropertyPlatformKey]
		platformMap, err := plaformKeyProp.getByteValue()
		if err != nil {
			return nil, err
		}
		idProp := arrP[playerStatsPropertyIdKey]
		id, err := idProp.getIntValue()
		if err != nil {
			return nil, err
		}
		isBotProp := arrP[playerStatsPropertybBotKey]
		isBot, err := isBotProp.getBoolValue()
		if err != nil {
			return nil, err
		}
		teamProp := arrP[playerStatsPropertyTeamKey]
		team, err := teamProp.getIntValue()
		if err != nil {
			return nil, err
		}
		assistsProp := arrP[playerStatsPropertyAssistsKey]
		assists, err := assistsProp.getIntValue()
		if err != nil {
			return nil, err
		}
		savesProp := arrP[playerStatsPropertySavesKey]
		saves, err := savesProp.getIntValue()
		if err != nil {
			return nil, err
		}
		shotsProp := arrP[playerStatsPropertyShotsKey]
		shots, err := shotsProp.getIntValue()
		if err != nil {
			return nil, err
		}
		goalsProp := arrP[playerStatsPropertyGoalsKey]
		goals, err := goalsProp.getIntValue()
		if err != nil {
			return nil, err
		}
		scoreProp := arrP[playerStatsPropertyScoreKey]
		score, err := scoreProp.getIntValue()
		if err != nil {
			return nil, err
		}
		stats = append(stats, PlayerStat{
			Name:     name,
			Plaform:  platformMap[playerStatsPropertyPlatformSubKey],
			OnlineId: id,
			IsBot:    isBot,
			Team:     team,
			Assists:  assists,
			Saves:    saves,
			Shots:    shots,
			Goals:    goals,
			Score:    score,
		})
	}

	return stats, nil
}
