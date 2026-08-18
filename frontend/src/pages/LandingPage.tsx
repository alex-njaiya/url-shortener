import FeaturesSection from "../components/shortener/FeaturesSection"
import Hero from "../components/shortener/Hero"
import ShortenForm from "../components/shortener/ShortenForm"

function LandingPage() {
  return (
    <div className="min-h-screen bg-paper">
      <Hero />
      <ShortenForm />
      <FeaturesSection />
    </div>
  )
}

export default LandingPage