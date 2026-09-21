// Web Audio API Emergency Siren Generator

class AudioAlertEngine {
  private ctx: AudioContext | null = null;
  private isPlaying = false;
  private sirenInterval: ReturnType<typeof setInterval> | null = null;
  private currentGain: GainNode | null = null;
  private masterVolume = 0.7;

  private getAudioContext(): AudioContext {
    if (!this.ctx) {
      const AudioCtx = window.AudioContext || (window as unknown as { webkitAudioContext: typeof AudioContext }).webkitAudioContext;
      this.ctx = new AudioCtx();
    }
    if (this.ctx.state === 'suspended') {
      this.ctx.resume().catch(() => {});
    }
    return this.ctx;
  }

  public setVolume(vol: number) {
    this.masterVolume = Math.max(0, Math.min(1, vol));
    if (this.currentGain && this.ctx) {
      this.currentGain.gain.setValueAtTime(this.masterVolume, this.ctx.currentTime);
    }
  }

  public getVolume(): number {
    return this.masterVolume;
  }

  /**
   * Starts playing continuous emergency siren until stop() is called.
   */
  public startSiren() {
    if (this.isPlaying) return;
    try {
      const ctx = this.getAudioContext();
      this.isPlaying = true;

      const gainNode = ctx.createGain();
      gainNode.gain.setValueAtTime(this.masterVolume, ctx.currentTime);
      gainNode.connect(ctx.destination);
      this.currentGain = gainNode;

      let highPitch = false;

      const playTonePair = () => {
        if (!this.isPlaying) return;
        const osc = ctx.createOscillator();
        const toneGain = ctx.createGain();

        const freq = highPitch ? 880 : 587.33; // A5 / D5 emergency siren pitch
        highPitch = !highPitch;

        osc.type = 'sawtooth';
        osc.frequency.setValueAtTime(freq, ctx.currentTime);

        // Fast pitch sweep down for dramatic emergency feel
        osc.frequency.exponentialRampToValueAtTime(freq * 0.85, ctx.currentTime + 0.35);

        toneGain.gain.setValueAtTime(0.3, ctx.currentTime);
        toneGain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.38);

        osc.connect(toneGain);
        toneGain.connect(gainNode);

        osc.start();
        osc.stop(ctx.currentTime + 0.4);
      };

      playTonePair();
      this.sirenInterval = setInterval(playTonePair, 420);
    } catch (e) {
      console.warn('Failed to start audio alert siren:', e);
    }
  }

  /**
   * Stops the currently playing siren immediately.
   */
  public stopSiren() {
    this.isPlaying = false;
    if (this.sirenInterval) {
      clearInterval(this.sirenInterval);
      this.sirenInterval = null;
    }
    if (this.currentGain && this.ctx) {
      try {
        this.currentGain.gain.linearRampToValueAtTime(0.001, this.ctx.currentTime + 0.1);
        setTimeout(() => {
          this.currentGain?.disconnect();
          this.currentGain = null;
        }, 150);
      } catch {
        this.currentGain = null;
      }
    }
  }

  /**
   * Plays a short test alert sound (2 pulses) to test audio & unlock browser AudioContext.
   */
  public playTestSound(): Promise<void> {
    return new Promise((resolve) => {
      this.startSiren();
      setTimeout(() => {
        this.stopSiren();
        resolve();
      }, 1200);
    });
  }
}

export const audioAlert = new AudioAlertEngine();
