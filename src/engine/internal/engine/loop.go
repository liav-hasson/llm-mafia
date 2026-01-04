package engine

// run serializes all state mutation and effect execution.
// It is the only place where GameState is modified.
// This is a two-phase executor:
//  1. Decision phase: command mutates state and returns effects
//  2. Effect phase: engine executes effects (Kafka, timers, etc.)
func (e *Engine) run() {
	for {
		select {
		case <-e.ctx.Done():
			return

		case cmd := <-e.cmdCh:
			// Phase 1: Apply command (pure state transformation)
			effects, err := cmd.Apply(e.state)
			if err != nil {
				// Command validation failed - do not execute effects
				// TODO: Add proper logging and error event emission
				_ = err
				continue
			}

			// Phase 2: Execute effects (side effects happen here)
			// side effects are any value that modifies an external system
			// (e.g. kafka publish) and or non-determenistic (e.g. timestamp)
			for _, effect := range effects {
				if err := effect.Execute(e.ctx, e.producer); err != nil {
					// Effect execution failed
					// TODO: Add retry logic, logging, metrics
					// Decision: continue with other effects or stop?
					_ = err
				}
			}

			// Phase 3: Schedule phase timer if phase changed
			// Cancel old timer and schedule new one based on current phase
			if _, isPhaseChange := cmd.(*PhaseChangeCommand); isPhaseChange {
				e.timers.CancelPhaseTimer()

				// Schedule timeout for the new phase (if applicable)
				timeout := GetPhaseTimeout(e.state.Phase)
				if timeout > 0 {
					nextPhase := GetNextPhase(e.state.Phase)
					e.timers.SchedulePhaseTimeout(
						e.state.Phase,
						e.state.Round,
						timeout,
						nextPhase,
						e.cmdCh,
						e.ctx,
					)
				}
			}

			// Also schedule timer when game starts
			if _, isStartGame := cmd.(*StartGameCommand); isStartGame {
				timeout := GetPhaseTimeout(e.state.Phase)
				if timeout > 0 {
					nextPhase := GetNextPhase(e.state.Phase)
					e.timers.SchedulePhaseTimeout(
						e.state.Phase,
						e.state.Round,
						timeout,
						nextPhase,
						e.cmdCh,
						e.ctx,
					)
				}
			}
		}
	}
}
