# BolusGPT

> [!CAUTION]
> Results are still a little buggy. Common issues are:
> - The GPT often doesn't want to call the API directly after web searching for nutrition info. You can remind it by just typing "api".
> - The GPT sometime ignores the instruction to never calculate a dose, and to solely rely on the API for that. Be vigilant that the result is coming from the actual API.

BolusGPT is an OpenAI custom [GPT](https://openai.com/index/introducing-gpts/) that calculates [bolus](https://en.wikipedia.org/wiki/Bolus_(medicine)) insulin doses via natural language. Users can use text/voice/images for prompts like the following:

> Bolus dose for an 8 oz steak, a cup of cooked broccoli, and a cup of brown rice.

> Change my insulin-to-carb ratio to 1:6.

The calculation is done on an HTTP server and the API is exposed to the GPT via an "Action" (function calling). The server integrates with Dexcom CGMs (via a port of [pydexcom](https://github.com/gagebenne/pydexcom/)) to get the user's real-time blood glucose level and trend, and also stores static settings like the user's insulin-to-carb ratio.

> [!NOTE]
> Insulin dosing is under the purview of the FDA, so users are required to self-host the server. All of the resources required to build the GPT, along with a [setup guide](./SETUP.md) is included in the repository. BolusGPT is not an FDA approved system, and is not sold or publicly hosted anywhere.

## Demo

TODO: Video

## Why?

I was diagnosed with Type 1 Diabetes, and found that manually counting carbs and calculating insulin doses multiple times a day sucks. I built this to make my life a little easier, and figured I would open source it in case it helps someone else.

## What is included?

### Server
1. Stores static settings (in a simple JSON file) required to calculate bolus doses .For example:
   1. Target blood glucose level
   1. Insulin-to-carb ratio
   1. Insulin sensitivity factor
1. Retrieves the user's real-time blood glucose level and trend via the Dexcom API.
1. Calculates bolus insulin doses via well-tested algorithm.

### GPT Resources
1. OpenAPI spec ([`openapi.yaml`](./openapi.yaml)) included so the GPT knows how to interact with the APIs.
1. Spec for the GPT itself ([`SPEC.md`](./SPEC.md))

## How does it work?

![](./bolusgpt.svg#gh-dark-mode-only)
![](./bolusgpt_light.svg#gh-light-mode-only)

A brief outline of how this works end-to-end:

1. The user follows the [setup guide](./SETUP.md) and creates their own BolusGPT in their OpenAI account.
1. On their preferred OpenAI app (iOS App, website, etc.), and using their preferred medium of interaction (text / voice), the user interacts with BolusGPT.
1. First, they onboard.
   1. BolusGPT asks for the user's insulin-to-carb ratio, target blood glucose level, etc.
   1. BolusGPT calls the `PATCH /me` API with the information, where it is stored in a JSON file on the server.
1. Next, they ask BolusGPT to dose their meal.
   1. BolusGPT asks what the user will eat, whether the user will soon be exercising, etc.
   1. BolusGPT searches the web for nutritional information to collect grams of carbs, fiber, protein, etc.
   1. BolusGPT calls the `POST /dose` API with the information.
      1. On the server, the current blood glucose level and trend is received from the Dexcom API.
      1. Using the nutrition info, exercise info, current blood glucose info, and stored user info, the insulin bolus is calculated and returned.
  1. BolusGPT presents the bolus dose to the user.
1. Optionally, the user can confirm they will use this dose (or tell BolusGPT they will opt for a different dose).
   1. If confirmed, BolusGPT can call `PATCH /me` with the dose used, and when. This can be an input to the next dose calculation (used to calculate insulin-on-board (IOB)).

### Documentation

#### `/me`

Any of the following fields can be retrieved (via `GET`), or updated (via `PATCH`):

- `fiber_multiplier` - Adjustment factor for dietary fiber's effect on insulin needs. A value of `1` counts all fiber. A value of `0` subtracts all fiber from total carbs
- `sugar_alcohol_multiplier` - Adjustment factor for sugar alcohols' impact on blood sugar. A value of `1` counts all sugar alcohol. A value of `0` subtracts all sugar alcohol from total carbs
- `protein_multiplier` - Factor representing how protein contributes to insulin demand. A value of `1` counts all protein. A value of `0` counts none of the protein
- `carb_threshold_to_count_protein_under` - Carb threshold under which protein is counted for dosing. For example, when the value is `20`, if the calculated carbs is under `20` protein is calculated according to the multiplier.
- `insulin_to_carb_ratio` - Grams of carbs covered by one unit of insulin. A value of `5` specifies 1 unit of insulin to 5 grams of carbs (1:5)
- `target_blood_glucose_level_in_mg_dl` - Target blood glucose level in mg/dL.
- `insulin_sensitivity_factor` - Blood glucose drop expected per unit of insulin. A value of `20` means a drop of 20 mg/dL is expected for 1 unit of insulin.
- `last_bolus_time` - Time of the last insulin bolus.
- `last_bolus_units_of_insulin` - Units of insulin used in the last bolus.

#### `/dose`

Any of the following fields can be provided when asking for a dose calculation via `POST`.

- `total_grams_of_carbs` - Total grams carbohydrates in the meal.
- `grams_of_fiber` - Grams of dietary fiber in the meal.
- `grams_of_sugar_alcohol` - Grams of sugar alcohols in the meal.
- `grams_of_protein` - Grams of protein in the meal.
- `minutes_of_exercise` - Duration of exercise in minutes that will occur after the bolus.
- `exercise_intensity` - Intensity of exercise that will occur after the bolus (`none`, `low`, `medium`, `high`).

### Why use OpenAI GPTs as an interface?

I wanted to make this quickly, and GPTs come with a lot for free, for example:

- A mobile app (the OpenAI app)
- Auth via OpenAI account
- Text, voice, images for modality
- Function calling with very little setup required
- Totally free

Someday I could see this turning into something more, but I am happy with where its at right now.

## How can I use this?

Please follow the [setup guide](./SETUP.md).

## Acknowledgements

1. Gage Benne for reverse engineering the Dexcom Share API: https://github.com/gagebenne/pydexcom
1. Gary Scheiner's book, [Think Like a Pancreas](https://www.amazon.com/Think-Like-Pancreas-Practical-Insulin-Completely/dp/0738215147)
1. Dr. Richard Bernstein's book, [Diabetes Solution](https://www.amazon.com/Dr-Bernsteins-Diabetes-Solution-Achieving/dp/0316182699)

## Future Directions

1. Use Open Food Facts Nutritional Data as a Knowledge Base: https://world.openfoodfacts.org/data
   1. Currently the file sizes are too large, even when chunked up.

## Todo

- Productionize
  - Put API on `api` subdomain
  - Direct `bolusgpt.com` to Github URL
