package enemies

import "crydes/helpers"

func (em *EnemiesManager) LoadAnimations() {
	em.loadSpiderAnimations()
	em.loadSkeletonAnimations()
	em.loadGoblinAnimations()
}

func (em *EnemiesManager) loadSpiderAnimations() {
	SPIRDER_idleRight := helpers.LoadAnimation("IDLE_R",
		"assets/spider/1.png",
		"assets/spider/2.png",
	)
	SPIRDER_moveRight := helpers.LoadAnimation("MOV_R",
		"assets/spider/9.png",
		"assets/spider/10.png",
		"assets/spider/11.png",
		"assets/spider/12.png",
	)
	SPIRDER_idleLeft := helpers.LoadAnimation("IDLE_L",
		"assets/spider/5.png",
		"assets/spider/6.png",
	)
	SPIRDER_moveLeft := helpers.LoadAnimation("MOV_L",
		"assets/spider/13.png",
		"assets/spider/14.png",
		"assets/spider/15.png",
		"assets/spider/16.png",
	)

	SPIDER_DEATH_LEFT := helpers.LoadAnimation("DEATH_L",
		"assets/spider/17.png",
		"assets/spider/18.png",
		"assets/spider/19.png",
		"assets/spider/20.png",
	)

	SPIDER_DEATH_RIGHT := helpers.LoadAnimation("DEATH_R",
		"assets/spider/21.png",
		"assets/spider/22.png",
		"assets/spider/23.png",
		"assets/spider/24.png",
	)

	em.Animations["spider"] = &map[string]*helpers.Animation{
		"idle_right":  SPIRDER_idleRight,
		"move_right":  SPIRDER_moveRight,
		"idle_left":   SPIRDER_idleLeft,
		"move_left":   SPIRDER_moveLeft,
		"death_left":  SPIDER_DEATH_LEFT,
		"death_right": SPIDER_DEATH_RIGHT,
	}

}

func (em *EnemiesManager) loadGoblinAnimations() {
	GOBLIN_idleRight := helpers.LoadAnimation("IDLE_R",
		"assets/goblin/1.png",
		"assets/goblin/2.png",
		"assets/goblin/3.png",
	)
	GOBLIN_moveRight := helpers.LoadAnimation("MOV_R",
		"assets/goblin/9.png",
		"assets/goblin/10.png",
		"assets/goblin/11.png",
		"assets/goblin/12.png",
	)
	GOBLIN_idleLeft := helpers.LoadAnimation("IDLE_L",
		"assets/goblin/5.png",
		"assets/goblin/6.png",
		"assets/goblin/7.png",
	)
	GOBLIN_moveLeft := helpers.LoadAnimation("MOV_L",
		"assets/goblin/13.png",
		"assets/goblin/14.png",
		"assets/goblin/15.png",
		"assets/goblin/16.png",
	)

	GOBLIN_DEATH_LEFT := helpers.LoadAnimation("DEATH_L",
		"assets/goblin/17.png",
		"assets/goblin/18.png",
		"assets/goblin/19.png",
		"assets/goblin/20.png",
	)

	GOBLIN_DEATH_RIGHT := helpers.LoadAnimation("DEATH_R",
		"assets/goblin/21.png",
		"assets/goblin/22.png",
		"assets/goblin/23.png",
		"assets/goblin/24.png",
	)

	em.Animations["goblin"] = &map[string]*helpers.Animation{
		"idle_right":  GOBLIN_idleRight,
		"move_right":  GOBLIN_moveRight,
		"idle_left":   GOBLIN_idleLeft,
		"move_left":   GOBLIN_moveLeft,
		"death_left":  GOBLIN_DEATH_LEFT,
		"death_right": GOBLIN_DEATH_RIGHT,
	}

}

func (em *EnemiesManager) loadSkeletonAnimations() {
	GOBLIN_idleRight := helpers.LoadAnimation("IDLE_R",
		"assets/skeleton/1.png",
		"assets/skeleton/2.png",
		"assets/skeleton/3.png",
	)
	GOBLIN_moveRight := helpers.LoadAnimation("MOV_R",
		"assets/skeleton/9.png",
		"assets/skeleton/10.png",
		"assets/skeleton/11.png",
		"assets/skeleton/12.png",
	)
	GOBLIN_idleLeft := helpers.LoadAnimation("IDLE_L",
		"assets/skeleton/5.png",
		"assets/skeleton/6.png",
		"assets/skeleton/7.png",
	)
	GOBLIN_moveLeft := helpers.LoadAnimation("MOV_L",
		"assets/skeleton/13.png",
		"assets/skeleton/14.png",
		"assets/skeleton/15.png",
		"assets/skeleton/16.png",
	)

	GOBLIN_DEATH_LEFT := helpers.LoadAnimation("DEATH_L",
		"assets/skeleton/17.png",
		"assets/skeleton/18.png",
		"assets/skeleton/19.png",
		"assets/skeleton/20.png",
	)

	GOBLIN_DEATH_RIGHT := helpers.LoadAnimation("DEATH_R",
		"assets/skeleton/21.png",
		"assets/skeleton/22.png",
		"assets/skeleton/23.png",
		"assets/skeleton/24.png",
	)

	em.Animations["skeleton"] = &map[string]*helpers.Animation{
		"idle_right":  GOBLIN_idleRight,
		"move_right":  GOBLIN_moveRight,
		"idle_left":   GOBLIN_idleLeft,
		"move_left":   GOBLIN_moveLeft,
		"death_left":  GOBLIN_DEATH_LEFT,
		"death_right": GOBLIN_DEATH_RIGHT,
	}

}
