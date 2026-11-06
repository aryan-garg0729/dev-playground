import { useSound } from '../stores/sound';
import { IconButton } from './Button';
import { SpeakerIcon, SpeakerOffIcon } from './icons';

/** Header mute switch for game sounds (on by default, remembered per device). */
export function SoundToggle() {
  const { muted, toggleMuted } = useSound();
  return (
    <IconButton aria-label={muted ? 'Unmute sounds' : 'Mute sounds'} aria-pressed={muted} onClick={toggleMuted}>
      {muted ? <SpeakerOffIcon /> : <SpeakerIcon />}
    </IconButton>
  );
}
